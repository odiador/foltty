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
  launch: (projectPath: string) => Promise<void>;
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
    launch: async (projectPath: string) => {
      await execFileAsync("open", [
        "-na",
        "Ghostty",
        "--args",
        `--working-directory=${projectPath}`,
      ]);
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
): Promise<void> {
  const selectedTerminal: Exclude<TerminalId, "auto"> =
    targetTerminal && targetTerminal !== "auto"
      ? targetTerminal
      : getActiveTerminalId();

  const definition = TERMINAL_REGISTRY[selectedTerminal];
  if (!definition) {
    throw new Error(`Unsupported terminal emulator: ${selectedTerminal}`);
  }

  await definition.launch(projectPath);
}
