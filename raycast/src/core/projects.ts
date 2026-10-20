import { execFile } from "child_process";
import fs from "fs";
import { homedir } from "os";
import path from "path";
import { promisify } from "util";

export const execFileAsync = promisify(execFile);

export const FOLTTY_FILE = path.join(homedir(), ".foltty_projects");
export const PMAN_FILE = path.join(homedir(), ".pman_projects");

// Determinar el archivo de persistencia principal
export function getStorageFilePath(): string {
  if (fs.existsSync(FOLTTY_FILE)) {
    return FOLTTY_FILE;
  }
  // Si existe .pman_projects pero no .foltty_projects, migramos/usamos pman
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

export interface FolttyProject {
  id: string;
  name: string;
  path: string;
  exists: boolean;
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
      const expandedPath = projPath.startsWith("~/")
        ? path.join(homedir(), projPath.slice(2))
        : projPath;

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
  const normalized = normalizePath(targetPath);

  if (!fs.existsSync(normalized)) {
    return { success: false, message: `La ruta no existe: ${normalized}` };
  }

  const stat = fs.statSync(normalized);
  if (!stat.isDirectory()) {
    return {
      success: false,
      message: `La ruta no es un directorio: ${normalized}`,
    };
  }

  const projects = getProjects();
  if (projects.some((p) => p.path === normalized)) {
    return {
      success: false,
      message: `El proyecto ya está registrado: ${normalized}`,
    };
  }

  const filePath = getStorageFilePath();
  fs.appendFileSync(filePath, `${normalized}\n`, "utf-8");

  // Mantener sincronizado pman si existe
  if (fs.existsSync(PMAN_FILE) && filePath !== PMAN_FILE) {
    try {
      fs.appendFileSync(PMAN_FILE, `${normalized}\n`, "utf-8");
    } catch {
      // Ignorar error secundario
    }
  }

  return {
    success: true,
    message: `Proyecto agregado: ${normalized}`,
    project: {
      id: normalized,
      name: path.basename(normalized) || normalized,
      path: normalized,
      exists: true,
    },
  };
}

/**
 * Elimina un proyecto por su ruta
 */
export function removeProject(targetPath: string): {
  success: boolean;
  message: string;
} {
  const normalized = normalizePath(targetPath);
  const projects = getProjects();
  const updated = projects.filter(
    (p) => p.path !== normalized && p.id !== targetPath,
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

  // Sincronizar pman si existe
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
  const normalized = normalizePath(projectPath);

  if (!fs.existsSync(normalized)) {
    throw new Error(`El directorio no existe en el disco: ${normalized}`);
  }

  await execFileAsync("open", [
    "-na",
    "Ghostty",
    "--args",
    `--working-directory=${normalized}`,
  ]);
}
