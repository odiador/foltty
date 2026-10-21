import fs from "fs";
import { homedir } from "os";
import path from "path";
import {
  openInTerminal,
  launchGhostty,
  getActiveTerminalId,
  type TerminalId,
  type GhosttyMode,
} from "./terminal";

export const FOLTTY_FILE = path.join(homedir(), ".foltty_projects");
export const PMAN_FILE = path.join(homedir(), ".pman_projects");
export const PRESETS_FILE = path.join(homedir(), ".foltty_presets.json");
export const NAMES_FILE = path.join(homedir(), ".foltty_names.json");
export const RECENT_FILE = path.join(homedir(), ".foltty_recent.json");

export type SortOrder = "recent" | "name-asc" | "name-desc" | "added";

export interface FolttyProject {
  id: string;
  name: string;
  customName?: string;
  path: string;
  exists: boolean;
  lastOpenedAt?: number;
}

export interface FolttyPreset {
  id: string;
  name: string;
  path: string;
  isCustom: boolean;
}

export function getStorageFilePath(): string {
  if (fs.existsSync(FOLTTY_FILE)) {
    return FOLTTY_FILE;
  }
  if (fs.existsSync(PMAN_FILE)) {
    try {
      const content = fs.readFileSync(PMAN_FILE, "utf-8");
      fs.writeFileSync(FOLTTY_FILE, content, "utf-8");
      return FOLTTY_FILE;
    } catch {
      return PMAN_FILE;
    }
  }
  return FOLTTY_FILE;
}

export function getCustomNames(): Record<string, string> {
  if (!fs.existsSync(NAMES_FILE)) {
    return {};
  }
  try {
    const data = fs.readFileSync(NAMES_FILE, "utf-8");
    return JSON.parse(data);
  } catch {
    return {};
  }
}

export function setProjectName(
  targetPath: string,
  customName: string,
): { success: boolean; message: string } {
  const normalized = normalizePath(targetPath);
  const names = getCustomNames();
  const trimmed = customName.trim();

  if (trimmed === "") {
    delete names[normalized];
  } else {
    names[normalized] = trimmed;
  }

  try {
    fs.writeFileSync(NAMES_FILE, JSON.stringify(names, null, 2), "utf-8");
    return {
      success: true,
      message: `Project renamed to '${trimmed || path.basename(normalized)}'`,
    };
  } catch (err) {
    return {
      success: false,
      message: `Error renaming project: ${String(err)}`,
    };
  }
}

export function getRecentHistory(): Record<string, number> {
  if (!fs.existsSync(RECENT_FILE)) {
    return {};
  }
  try {
    const data = fs.readFileSync(RECENT_FILE, "utf-8");
    return JSON.parse(data);
  } catch {
    return {};
  }
}

export function recordProjectOpened(targetPath: string): void {
  try {
    const normalized = normalizePath(targetPath);
    const history = getRecentHistory();
    history[normalized] = Date.now();
    fs.writeFileSync(RECENT_FILE, JSON.stringify(history, null, 2), "utf-8");
  } catch (err) {
    console.error("Error updating project recent timestamp:", err);
  }
}

export function formatRelativeTime(timestamp?: number): string | undefined {
  if (!timestamp || typeof timestamp !== "number") {
    return undefined;
  }
  const now = Date.now();
  const diffMs = now - timestamp;
  if (diffMs < 0) {
    return "Just now";
  }

  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHour = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHour / 24);

  if (diffMin < 1) {
    return "Just now";
  }
  if (diffMin < 60) {
    return `${diffMin}m ago`;
  }
  if (diffHour < 24) {
    return `${diffHour}h ago`;
  }
  if (diffDay === 1) {
    return "Yesterday";
  }
  if (diffDay < 7) {
    return `${diffDay}d ago`;
  }

  const d = new Date(timestamp);
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric" });
}

export function sortProjects(
  projects: FolttyProject[],
  sortOrder: SortOrder = "recent",
): FolttyProject[] {
  const copy = [...projects];
  switch (sortOrder) {
    case "recent":
      return copy.sort((a, b) => {
        const timeA = a.lastOpenedAt || 0;
        const timeB = b.lastOpenedAt || 0;
        if (timeB !== timeA) {
          return timeB - timeA;
        }
        return a.name.localeCompare(b.name, undefined, { sensitivity: "base" });
      });
    case "name-asc":
      return copy.sort((a, b) =>
        a.name.localeCompare(b.name, undefined, { sensitivity: "base" }),
      );
    case "name-desc":
      return copy.sort((a, b) =>
        b.name.localeCompare(a.name, undefined, { sensitivity: "base" }),
      );
    case "added":
      return copy;
    default:
      return copy;
  }
}

export function getDefaultPresets(): FolttyPreset[] {
  const home = homedir();
  const presets: FolttyPreset[] = [
    { id: "home", name: "Home (~)", path: home, isCustom: false },
    {
      id: "documents",
      name: "Documents",
      path: path.join(home, "Documents"),
      isCustom: false,
    },
    {
      id: "downloads",
      name: "Downloads",
      path: path.join(home, "Downloads"),
      isCustom: false,
    },
    {
      id: "desktop",
      name: "Desktop",
      path: path.join(home, "Desktop"),
      isCustom: false,
    },
    {
      id: "config",
      name: "Config (.config)",
      path: path.join(home, ".config"),
      isCustom: false,
    },
  ];

  const devDirs = ["Development", "Developer", "Projects", "Code", "workspace"];
  for (const dir of devDirs) {
    const fullPath = path.join(home, dir);
    if (fs.existsSync(fullPath) && fs.statSync(fullPath).isDirectory()) {
      presets.push({
        id: "dev",
        name: `Development (${dir})`,
        path: fullPath,
        isCustom: false,
      });
      break;
    }
  }

  return presets;
}

export function getCustomPresets(): Record<string, string> {
  if (!fs.existsSync(PRESETS_FILE)) {
    return {};
  }
  try {
    const data = fs.readFileSync(PRESETS_FILE, "utf-8");
    return JSON.parse(data);
  } catch {
    return {};
  }
}

export function getAllPresets(): FolttyPreset[] {
  const defaults = getDefaultPresets();
  const customs = getCustomPresets();

  const customPresets: FolttyPreset[] = Object.entries(customs).map(
    ([id, p]) => ({
      id,
      name: `${id} (Custom)`,
      path: normalizePath(p),
      isCustom: true,
    }),
  );

  return [...defaults, ...customPresets];
}

export function addCustomPreset(
  name: string,
  targetPath: string,
): { success: boolean; message: string } {
  const trimmed = name.trim().toLowerCase();
  if (!trimmed) {
    return { success: false, message: "Preset name cannot be empty" };
  }

  const normalized = normalizePath(targetPath);
  if (!fs.existsSync(normalized)) {
    return { success: false, message: `Path does not exist: ${normalized}` };
  }

  const customs = getCustomPresets();
  customs[trimmed] = normalized;

  try {
    fs.writeFileSync(PRESETS_FILE, JSON.stringify(customs, null, 2), "utf-8");
    return {
      success: true,
      message: `Preset saved: ${trimmed} -> ${normalized}`,
    };
  } catch (err) {
    return { success: false, message: `Error saving preset: ${String(err)}` };
  }
}

export function removeCustomPreset(name: string): {
  success: boolean;
  message: string;
} {
  const trimmed = name.trim().toLowerCase();
  const customs = getCustomPresets();

  if (!(trimmed in customs)) {
    return { success: false, message: `Preset not found: ${name}` };
  }

  delete customs[trimmed];
  try {
    fs.writeFileSync(PRESETS_FILE, JSON.stringify(customs, null, 2), "utf-8");
    return { success: true, message: `Preset removed: ${trimmed}` };
  } catch (err) {
    return { success: false, message: `Error removing preset: ${String(err)}` };
  }
}

export function normalizePath(inputPath: string): string {
  let resolved = inputPath.trim();
  if (resolved.startsWith("~/")) {
    resolved = path.join(homedir(), resolved.slice(2));
  } else if (resolved === "~") {
    resolved = homedir();
  } else {
    resolved = path.resolve(resolved);
  }
  return resolved;
}

export function resolvePresetOrPath(input: string): string {
  const trimmed = input.trim();
  const lower = trimmed.toLowerCase();

  const presets = getAllPresets();
  for (const p of presets) {
    if (
      p.id.toLowerCase() === lower ||
      (lower === "docs" && p.id === "documents") ||
      (lower === "dl" && p.id === "downloads")
    ) {
      return p.path;
    }
  }

  return normalizePath(trimmed);
}

export function getProjects(): FolttyProject[] {
  const filePath = getStorageFilePath();
  if (!fs.existsSync(filePath)) {
    return [];
  }

  try {
    const content = fs.readFileSync(filePath, "utf-8");
    const customNames = getCustomNames();
    const recentHistory = getRecentHistory();

    const lines = content
      .split("\n")
      .map((line) => line.trim())
      .filter((line) => line.length > 0 && !line.startsWith("#"));

    const uniquePaths = Array.from(new Set(lines));

    return uniquePaths.map((projPath) => {
      const expandedPath = normalizePath(projPath);
      const exists = fs.existsSync(expandedPath);
      let name = path.basename(expandedPath) || expandedPath;
      let customLabel: string | undefined = undefined;

      if (customNames[expandedPath]) {
        name = customNames[expandedPath];
        customLabel = customNames[expandedPath];
      }

      const lastOpenedAt = recentHistory[expandedPath];

      return {
        id: expandedPath,
        name,
        customName: customLabel,
        path: expandedPath,
        exists,
        lastOpenedAt,
      };
    });
  } catch (error) {
    console.error("Error reading foltty projects:", error);
    return [];
  }
}

export function addProject(
  targetPath: string,
  customName?: string,
): { success: boolean; message: string; project?: FolttyProject } {
  const resolved = resolvePresetOrPath(targetPath);

  if (!fs.existsSync(resolved)) {
    return { success: false, message: `Path does not exist: ${resolved}` };
  }

  const stat = fs.statSync(resolved);
  if (!stat.isDirectory()) {
    return { success: false, message: `Path is not a directory: ${resolved}` };
  }

  const projects = getProjects();
  if (!projects.some((p) => p.path === resolved)) {
    const filePath = getStorageFilePath();
    fs.appendFileSync(filePath, `${resolved}\n`, "utf-8");

    if (fs.existsSync(PMAN_FILE) && filePath !== PMAN_FILE) {
      try {
        fs.appendFileSync(PMAN_FILE, `${resolved}\n`, "utf-8");
      } catch {
        // Ignore secondary error
      }
    }
  }

  if (customName && customName.trim() !== "") {
    setProjectName(resolved, customName);
  }

  return {
    success: true,
    message: `Project added: ${resolved}`,
    project: {
      id: resolved,
      name: customName || path.basename(resolved) || resolved,
      customName,
      path: resolved,
      exists: true,
    },
  };
}

export function removeProject(targetPath: string): {
  success: boolean;
  message: string;
} {
  const resolved = resolvePresetOrPath(targetPath);
  const projects = getProjects();
  const updated = projects.filter(
    (p) => p.path !== resolved && p.name !== targetPath && p.id !== targetPath,
  );

  if (updated.length === projects.length) {
    return {
      success: false,
      message: `Project not found in list: ${targetPath}`,
    };
  }

  setProjectName(resolved, ""); // clean up custom name

  try {
    const history = getRecentHistory();
    if (history[resolved]) {
      delete history[resolved];
      fs.writeFileSync(RECENT_FILE, JSON.stringify(history, null, 2), "utf-8");
    }
  } catch {
    // Ignore cleanup error
  }

  const filePath = getStorageFilePath();
  const content =
    updated.map((p) => p.path).join("\n") + (updated.length > 0 ? "\n" : "");
  fs.writeFileSync(filePath, content, "utf-8");

  if (fs.existsSync(PMAN_FILE) && filePath !== PMAN_FILE) {
    try {
      fs.writeFileSync(PMAN_FILE, content, "utf-8");
    } catch {
      // Ignore secondary error
    }
  }

  return { success: true, message: `Project removed successfully.` };
}

export async function openProjectInTerminal(
  projectPath: string,
  terminal?: TerminalId,
  ghosttyMode?: GhosttyMode,
  title?: string,
): Promise<void> {
  const resolved = resolvePresetOrPath(projectPath);

  if (!fs.existsSync(resolved)) {
    throw new Error(`Directory does not exist on disk: ${resolved}`);
  }

  recordProjectOpened(resolved);

  const selectedTerminal =
    terminal && terminal !== "auto" ? terminal : getActiveTerminalId();

  if (selectedTerminal === "ghostty") {
    await launchGhostty(resolved, ghosttyMode, title);
    return;
  }

  await openInTerminal(resolved, selectedTerminal, title);
}

export async function openInGhostty(
  projectPath: string,
  mode?: GhosttyMode,
  title?: string,
): Promise<void> {
  await openProjectInTerminal(projectPath, "ghostty", mode, title);
}
