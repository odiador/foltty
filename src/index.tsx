import {
  ActionPanel,
  Action,
  List,
  Icon,
  Color,
  showToast,
  Toast,
  closeMainWindow,
  useNavigation,
  Form,
} from "@raycast/api";
import { execFile } from "child_process";
import fs from "fs";
import { homedir } from "os";
import path from "path";
import { useState, useEffect, useCallback } from "react";
import { promisify } from "util";

const execFileAsync = promisify(execFile);
const PMAN_FILE = path.join(homedir(), ".pman_projects");

interface Project {
  id: string;
  name: string;
  path: string;
  exists: boolean;
}

// Lee los proyectos de ~/.pman_projects
function readProjectsFromFile(): Project[] {
  if (!fs.existsSync(PMAN_FILE)) {
    return [];
  }

  try {
    const content = fs.readFileSync(PMAN_FILE, "utf-8");
    const lines = content
      .split("\n")
      .map((line) => line.trim())
      .filter((line) => line.length > 0 && !line.startsWith("#"));

    // Eliminar duplicados manteniendo el orden
    const uniquePaths = Array.from(new Set(lines));

    return uniquePaths.map((projPath) => {
      // Normalizar ~ a homedir si es necesario
      const expandedPath = projPath.startsWith("~/")
        ? path.join(homedir(), projPath.slice(2))
        : projPath;

      const exists = fs.existsSync(expandedPath);
      return {
        id: projPath,
        name: path.basename(expandedPath) || expandedPath,
        path: expandedPath,
        exists,
      };
    });
  } catch (error) {
    console.error("Error al leer archivo pman:", error);
    return [];
  }
}

// Guarda la lista de proyectos en ~/.pman_projects
function saveProjectsToFile(projects: Project[]): void {
  const content =
    projects.map((p) => p.path).join("\n") + (projects.length > 0 ? "\n" : "");
  fs.writeFileSync(PMAN_FILE, content, "utf-8");
}

export default function Command() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadProjects = useCallback(() => {
    setIsLoading(true);
    const list = readProjectsFromFile();
    setProjects(list);
    setIsLoading(false);
  }, []);

  useEffect(() => {
    loadProjects();
  }, [loadProjects]);

  // Abrir la ruta seleccionada en Ghostty con ventana limpia
  const openInGhostty = async (projectPath: string) => {
    if (!fs.existsSync(projectPath)) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Directorio no encontrado",
        message: projectPath,
      });
      return;
    }

    try {
      await closeMainWindow();
      await execFileAsync("open", [
        "-na",
        "Ghostty",
        "--args",
        `--working-directory=${projectPath}`,
      ]);
    } catch (error) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Error al abrir Ghostty",
        message: error instanceof Error ? error.message : String(error),
      });
    }
  };

  // Eliminar el proyecto de la lista y del archivo
  const deleteProject = async (project: Project) => {
    const updated = projects.filter((p) => p.id !== project.id);
    setProjects(updated);
    saveProjectsToFile(updated);

    await showToast({
      style: Toast.Style.Success,
      title: "Proyecto removido",
      message: project.name,
    });
  };

  return (
    <List
      isLoading={isLoading}
      searchBarPlaceholder="Buscar o filtrar proyectos por nombre o ruta..."
    >
      {projects.length === 0 && !isLoading ? (
        <List.EmptyView
          icon={Icon.Folder}
          title="No hay proyectos registrados"
          description="Presiona ⌘+N para agregar un proyecto, o usa 'padd' en tu terminal."
          actions={
            <ActionPanel>
              <Action.Push
                title="Agregar Proyecto"
                icon={Icon.Plus}
                target={<AddProjectForm onProjectAdded={loadProjects} />}
              />
              <Action.OpenWith
                title="Abrir ~/.pman_projects"
                path={PMAN_FILE}
                shortcut={{ modifiers: ["cmd", "shift"], key: "p" }}
              />
            </ActionPanel>
          }
        />
      ) : (
        projects.map((project) => (
          <List.Item
            key={project.id}
            icon={
              project.exists
                ? Icon.Folder
                : { source: Icon.ExclamationMark, tintColor: Color.Red }
            }
            title={project.name}
            subtitle={project.path.replace(homedir(), "~")}
            accessories={[
              !project.exists
                ? {
                    tag: { value: "Ruta inaccesible", color: Color.Red },
                    tooltip: "La carpeta no existe en el disco",
                  }
                : {},
            ]}
            actions={
              <ActionPanel>
                <ActionPanel.Section title="Ghostty">
                  <Action
                    title="Abrir En Ghostty"
                    icon={Icon.Terminal}
                    onAction={() => openInGhostty(project.path)}
                  />
                </ActionPanel.Section>

                <ActionPanel.Section title="Acciones de Archivo">
                  <Action.ShowInFinder
                    path={project.path}
                    shortcut={{ modifiers: ["cmd"], key: "enter" }}
                  />
                  <Action.OpenWith
                    path={project.path}
                    shortcut={{ modifiers: ["cmd", "shift"], key: "o" }}
                  />
                  <Action.CopyToClipboard
                    title="Copiar Ruta"
                    content={project.path}
                    shortcut={{ modifiers: ["cmd"], key: "c" }}
                  />
                </ActionPanel.Section>

                <ActionPanel.Section title="Gestión de Proyectos">
                  <Action.Push
                    title="Agregar Proyecto"
                    icon={Icon.Plus}
                    shortcut={{ modifiers: ["cmd"], key: "n" }}
                    target={<AddProjectForm onProjectAdded={loadProjects} />}
                  />
                  <Action
                    title="Remover De La Lista"
                    icon={Icon.Trash}
                    style={Action.Style.Destructive}
                    shortcut={{ modifiers: ["cmd"], key: "backspace" }}
                    onAction={() => deleteProject(project)}
                  />
                  <Action
                    title="Recargar Lista"
                    icon={Icon.ArrowClockwise}
                    shortcut={{ modifiers: ["cmd"], key: "r" }}
                    onAction={loadProjects}
                  />
                </ActionPanel.Section>
              </ActionPanel>
            }
          />
        ))
      )}
    </List>
  );
}

function AddProjectForm({ onProjectAdded }: { onProjectAdded: () => void }) {
  const { pop } = useNavigation();
  const [folderPath, setFolderPath] = useState<string[]>([]);
  const [customPath, setCustomPath] = useState<string>("");

  const handleSubmit = async () => {
    let target = folderPath.length > 0 ? folderPath[0] : customPath.trim();

    if (!target) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Por favor selecciona o escribe la ruta de un proyecto",
      });
      return;
    }

    if (target.startsWith("~/")) {
      target = path.join(homedir(), target.slice(2));
    }

    if (!fs.existsSync(target)) {
      await showToast({
        style: Toast.Style.Failure,
        title: "La carpeta no existe",
        message: target,
      });
      return;
    }

    const currentProjects = readProjectsFromFile();
    if (currentProjects.some((p) => p.path === target)) {
      await showToast({
        style: Toast.Style.Failure,
        title: "El proyecto ya está en la lista",
      });
      return;
    }

    // Append to file
    fs.appendFileSync(PMAN_FILE, `${target}\n`, "utf-8");

    await showToast({
      style: Toast.Style.Success,
      title: "Proyecto agregado exitosamente",
      message: path.basename(target),
    });

    onProjectAdded();
    pop();
  };

  return (
    <Form
      actions={
        <ActionPanel>
          <Action.SubmitForm title="Guardar Proyecto" onSubmit={handleSubmit} />
        </ActionPanel>
      }
    >
      <Form.Description text="Selecciona una carpeta usando el explorador de archivos o ingresa una ruta manual." />
      <Form.FilePicker
        id="folderPicker"
        title="Seleccionar Carpeta"
        canChooseDirectories={true}
        canChooseFiles={false}
        allowMultipleSelection={false}
        value={folderPath}
        onChange={setFolderPath}
      />
      <Form.TextField
        id="customPath"
        title="O escribir Ruta"
        placeholder="/Users/usuario/Proyectos/mi-proyecto"
        value={customPath}
        onChange={setCustomPath}
      />
    </Form>
  );
}
