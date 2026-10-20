import { getPreferenceValues } from "@raycast/api";
import { execFile } from "child_process";
import fs from "fs";
import { promisify } from "util";

const execFileAsync = promisify(execFile);

export type TerminalId =
  | "auto"
  | "ghostty"
  | "iterm"
  | "wezterm"
  | "alacritty"
  | "kitty"
  | "warp"
  | "terminal";

export interface TerminalDefinition {
  id: TerminalId;
  name: string;
  appPath: string;
  isAvailable: () => boolean;
  launch: (projectPath: string, title?: string) => Promise<void>;
}

export function getGhosttyTargetPreference(): "tab" | "window" {
  try {
    const prefs = getPreferenceValues<{ ghosttyTarget?: "tab" | "window" }>();
    if (prefs.ghosttyTarget === "tab") {
      return "tab";
    }
    if (prefs.ghosttyTarget === "window") {
      return "window";
    }
  } catch {
    // fallback
  }
  return "window";
}

export async function launchGhostty(
  projectPath: string,
  mode?: "tab" | "window",
  title?: string,
): Promise<void> {
  const targetMode = mode || getGhosttyTargetPreference();
  const targetTitle = title || "";
  const script = `
on run argv
    set targetDir to item 1 of argv
    set openMode to item 2 of argv
    set targetTitle to ""
    if (count of argv) >= 3 then
        set targetTitle to item 3 of argv
    end if

    tell application "Ghostty"
        activate
        set cfg to new surface configuration
        set initial working directory of cfg to targetDir
        if openMode is "tab" and (count of windows) > 0 then
            set myTab to new tab in front window with configuration cfg
            set myTerm to focused terminal of myTab
        else
            set myWin to new window with configuration cfg
            set myTab to selected tab of myWin
            set myTerm to focused terminal of myTab
        end if
        if targetTitle is not "" then
            perform action ("set_tab_title:" & targetTitle) on myTerm
        end if
    end tell
end run`;

  try {
    await execFileAsync("osascript", [
      "-e",
      script,
      projectPath,
      targetMode,
      targetTitle,
    ]);
  } catch {
    // Fallback to macOS open without -n to preserve instance
    await execFileAsync("open", ["-a", "Ghostty", projectPath]);
  }
}

export const TERMINAL_REGISTRY: Record<
  Exclude<TerminalId, "auto">,
  TerminalDefinition
> = {
  ghostty: {
    id: "ghostty",
    name: "Ghostty",
    appPath: "/Applications/Ghostty.app",
    isAvailable: () =>
      fs.existsSync("/Applications/Ghostty.app") ||
      fs.existsSync(`${process.env.HOME}/Applications/Ghostty.app`),
    launch: async (projectPath: string, title?: string) => {
      await launchGhostty(projectPath, undefined, title);
    },
  },
  iterm: {
    id: "iterm",
    name: "iTerm2",
    appPath: "/Applications/iTerm.app",
    isAvailable: () =>
      fs.existsSync("/Applications/iTerm.app") ||
      fs.existsSync(`${process.env.HOME}/Applications/iTerm.app`),
    launch: async (projectPath: string) => {
      await execFileAsync("open", ["-a", "iTerm", projectPath]);
    },
  },
  wezterm: {
    id: "wezterm",
    name: "WezTerm",
    appPath: "/Applications/WezTerm.app",
    isAvailable: () =>
      fs.existsSync("/Applications/WezTerm.app") ||
      fs.existsSync(`${process.env.HOME}/Applications/WezTerm.app`),
    launch: async (projectPath: string) => {
      await execFileAsync("open", [
        "-na",
        "WezTerm",
        "--args",
        "start",
        "--cwd",
        projectPath,
      ]);
    },
  },
  alacritty: {
    id: "alacritty",
    name: "Alacritty",
    appPath: "/Applications/Alacritty.app",
    isAvailable: () =>
      fs.existsSync("/Applications/Alacritty.app") ||
      fs.existsSync(`${process.env.HOME}/Applications/Alacritty.app`),
    launch: async (projectPath: string) => {
      await execFileAsync("open", [
        "-na",
        "Alacritty",
        "--args",
        "--working-directory",
        projectPath,
      ]);
    },
  },
  kitty: {
    id: "kitty",
    name: "Kitty",
    appPath: "/Applications/kitty.app",
    isAvailable: () =>
      fs.existsSync("/Applications/kitty.app") ||
      fs.existsSync(`${process.env.HOME}/Applications/kitty.app`),
    launch: async (projectPath: string) => {
      await execFileAsync("open", [
        "-na",
        "kitty",
        "--args",
        "--directory",
        projectPath,
      ]);
    },
  },
  warp: {
    id: "warp",
    name: "Warp",
    appPath: "/Applications/Warp.app",
    isAvailable: () =>
      fs.existsSync("/Applications/Warp.app") ||
      fs.existsSync(`${process.env.HOME}/Applications/Warp.app`),
    launch: async (projectPath: string) => {
      await execFileAsync("open", ["-na", "Warp", projectPath]);
    },
  },
  terminal: {
    id: "terminal",
    name: "Apple Terminal",
    appPath: "/System/Applications/Utilities/Terminal.app",
    isAvailable: () => true,
    launch: async (projectPath: string) => {
      await execFileAsync("open", ["-a", "Terminal", projectPath]);
    },
  },
};

interface Preferences {
  defaultTerminal?: TerminalId;
}

export function detectDefaultTerminal(): Exclude<TerminalId, "auto"> {
  const priorityOrder: Exclude<TerminalId, "auto">[] = [
    "ghostty",
    "iterm",
    "wezterm",
    "warp",
    "alacritty",
    "kitty",
    "terminal",
  ];

  for (const termId of priorityOrder) {
    if (TERMINAL_REGISTRY[termId].isAvailable()) {
      return termId;
    }
  }

  return "terminal";
}

export function getActiveTerminalId(): Exclude<TerminalId, "auto"> {
  try {
    const prefs = getPreferenceValues<Preferences>();
    if (prefs.defaultTerminal && prefs.defaultTerminal !== "auto") {
      return prefs.defaultTerminal;
    }
  } catch {
    // If running in context without preferences
  }

  return detectDefaultTerminal();
}

export function getActiveTerminalName(): string {
  const termId = getActiveTerminalId();
  return TERMINAL_REGISTRY[termId].name;
}

export async function openInTerminal(
  projectPath: string,
  targetTerminal?: TerminalId,
  title?: string,
): Promise<void> {
  const selectedTerminal: Exclude<TerminalId, "auto"> =
    targetTerminal && targetTerminal !== "auto"
      ? targetTerminal
      : getActiveTerminalId();

  const definition = TERMINAL_REGISTRY[selectedTerminal];
  if (!definition) {
    throw new Error(`Unsupported terminal emulator: ${selectedTerminal}`);
  }

  await definition.launch(projectPath, title);
}
