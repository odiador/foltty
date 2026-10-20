import { execFile } from "child_process";
import fs from "fs";
import { homedir } from "os";
import path from "path";
import { promisify } from "util";

export const execFileAsync = promisify(execFile);

export const FOLTTY_FILE = path.join(homedir(), ".foltty_projects");
export const PMAN_FILE = path.join(homedir(), ".pman_projects");
export const PRESETS_FILE = path.join(homedir(), ".foltty_presets.json");

export interface FolttyProject {
  id: string;
  name: string;
  path: string;
  exists: boolean;
}

export interface FolttyPreset {
  id: string;
  name: string;
  path: string;
  isCustom: boolean;
}

// Determinar el archivo de persistencia principal
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

export function getDefaultPresets(): FolttyPreset[] {
  const home = homedir();
  const presets: FolttyPreset[] = [
    { id: "home", name: "Home (~)", path: home, isCustom: false },
    {
      id: "documents",
      name: "Documentos",
      path: path.join(home, "Documents"),
      isCustom: false,
    },
    {
      id: "downloads",
      name: "Descargas",
      path: path.join(home, "Downloads"),
      isCustom: false,
    },
    {
      id: "desktop",
      name: "Escritorio",
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
        name: `Desarrollo (${dir})`,
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
    return {
      success: false,
      message: "El nombre del preset no puede estar vacío",
    };
  }

  const normalized = normalizePath(targetPath);
  if (!fs.existsSync(normalized)) {
    return { success: false, message: `La ruta no existe: ${normalized}` };
  }

  const customs = getCustomPresets();
  customs[trimmed] = normalized;

  try {
    fs.writeFileSync(PRESETS_FILE, JSON.stringify(customs, null, 2), "utf-8");
    return {
      success: true,
      message: `Preset guardado: ${trimmed} -> ${normalized}`,
    };
  } catch (err) {
    return {
      success: false,
      message: `Error al guardar preset: ${String(err)}`,
    };
  }
}

export function removeCustomPreset(name: string): {
  success: boolean;
  message: string;
} {
  const trimmed = name.trim().toLowerCase();
  const customs = getCustomPresets();

  if (!(trimmed in customs)) {
    return { success: false, message: `Preset no encontrado: ${name}` };
  }

  delete customs[trimmed];
  try {
    fs.writeFileSync(PRESETS_FILE, JSON.stringify(customs, null, 2), "utf-8");
    return { success: true, message: `Preset removido: ${trimmed}` };
  } catch (err) {
    return {
      success: false,
      message: `Error al remover preset: ${String(err)}`,
    };
  }
}

/**
 * Normaliza y expande rutas con ~
 */
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

/**
 * Obtiene la lista ordenada de proyectos registrados
 */
export function getProjects(): FolttyProject[] {
  const filePath = getStorageFilePath();
  if (!fs.existsSync(filePath)) {
    return [];
  }

  try {
    const content = fs.readFileSync(filePath, "utf-8");
    const lines = content
      .split("\n")
      .map((line) => line.trim())
      .filter((line) => line.length > 0 && !line.startsWith("#"));

    const uniquePaths = Array.from(new Set(lines));

    return uniquePaths.map((projPath) => {
      const expandedPath = normalizePath(projPath);
      const exists = fs.existsSync(expandedPath);
      return {
        id: expandedPath,
        name: path.basename(expandedPath) || expandedPath,
        path: expandedPath,
        exists,
      };
    });
  } catch (error) {
    console.error("Error al leer proyectos de foltty:", error);
    return [];
  }
}

/**
 * Agrega un proyecto a la lista si no existe
 */
export function addProject(targetPath: string): {
  success: boolean;
  message: string;
  project?: FolttyProject;
} {
  const resolved = resolvePresetOrPath(targetPath);

  if (!fs.existsSync(resolved)) {
    return { success: false, message: `La ruta no existe: ${resolved}` };
  }

  const stat = fs.statSync(resolved);
  if (!stat.isDirectory()) {
    return {
      success: false,
      message: `La ruta no es un directorio: ${resolved}`,
    };
  }

  const projects = getProjects();
  if (projects.some((p) => p.path === resolved)) {
    return {
      success: false,
      message: `El proyecto ya está registrado: ${resolved}`,
    };
  }

  const filePath = getStorageFilePath();
  fs.appendFileSync(filePath, `${resolved}\n`, "utf-8");

  if (fs.existsSync(PMAN_FILE) && filePath !== PMAN_FILE) {
    try {
      fs.appendFileSync(PMAN_FILE, `${resolved}\n`, "utf-8");
    } catch {
      // Ignorar error secundario
    }
  }

  return {
    success: true,
    message: `Proyecto agregado: ${resolved}`,
    project: {
      id: resolved,
      name: path.basename(resolved) || resolved,
      path: resolved,
      exists: true,
    },
  };
}

/**
 * Elimina un proyecto por su ruta o nombre
 */
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
      message: `Proyecto no encontrado en la lista: ${targetPath}`,
    };
  }

  const filePath = getStorageFilePath();
  const content =
    updated.map((p) => p.path).join("\n") + (updated.length > 0 ? "\n" : "");
  fs.writeFileSync(filePath, content, "utf-8");

  if (fs.existsSync(PMAN_FILE) && filePath !== PMAN_FILE) {
    try {
      fs.writeFileSync(PMAN_FILE, content, "utf-8");
    } catch {
      // Ignorar error secundario
    }
  }

  return { success: true, message: `Proyecto removido exitosamente.` };
}

/**
 * Abre la ruta especificada en una instancia limpia de Ghostty
 */
export async function openInGhostty(projectPath: string): Promise<void> {
  const resolved = resolvePresetOrPath(projectPath);

  if (!fs.existsSync(resolved)) {
    throw new Error(`El directorio no existe en el disco: ${resolved}`);
  }

  await execFileAsync("open", [
    "-na",
    "Ghostty",
    "--args",
    `--working-directory=${resolved}`,
  ]);
}
