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
import { homedir } from "os";
import path from "path";
import { useState, useEffect, useCallback } from "react";
import {
  getProjects,
  addProject,
  removeProject,
  openInGhostty,
  getStorageFilePath,
  type FolttyProject,
} from "./core/projects";

export default function Command() {
  const [projects, setProjects] = useState<FolttyProject[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadProjects = useCallback(() => {
    setIsLoading(true);
    const list = getProjects();
    setProjects(list);
    setIsLoading(false);
  }, []);

  useEffect(() => {
    loadProjects();
  }, [loadProjects]);

  const handleOpenGhostty = async (projectPath: string) => {
    try {
      await closeMainWindow();
      await openInGhostty(projectPath);
    } catch (error) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Error al abrir Ghostty",
        message: error instanceof Error ? error.message : String(error),
      });
    }
  };

  const handleDelete = async (project: FolttyProject) => {
    const res = removeProject(project.path);
    if (res.success) {
      loadProjects();
      await showToast({
        style: Toast.Style.Success,
        title: "Proyecto Removido",
        message: project.name,
      });
    } else {
      await showToast({
        style: Toast.Style.Failure,
        title: "Error al Remover",
        message: res.message,
      });
    }
  };

  const storageFile = getStorageFilePath();

  return (
    <List
      isLoading={isLoading}
      searchBarPlaceholder="Buscar o filtrar proyectos por nombre o ruta..."
    >
      {projects.length === 0 && !isLoading ? (
        <List.EmptyView
          icon={Icon.Folder}
          title="No hay proyectos registrados en Foltty"
          description="Presiona ⌘+N para agregar un proyecto, o usa 'foltty add' en tu terminal."
          actions={
            <ActionPanel>
              <Action.Push
                title="Agregar Proyecto"
                icon={Icon.Plus}
                target={<AddProjectForm onProjectAdded={loadProjects} />}
              />
              <Action.OpenWith
                title="Abrir Archivo De Proyectos"
                path={storageFile}
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
                    onAction={() => handleOpenGhostty(project.path)}
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
                    onAction={() => handleDelete(project)}
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
    const target = folderPath.length > 0 ? folderPath[0] : customPath.trim();

    if (!target) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Por favor selecciona o escribe la ruta de un proyecto",
      });
      return;
    }

    const res = addProject(target);
    if (!res.success) {
      await showToast({
        style: Toast.Style.Failure,
        title: "No se pudo agregar el proyecto",
        message: res.message,
      });
      return;
    }

    await showToast({
      style: Toast.Style.Success,
      title: "Proyecto Agregado",
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
        title="O Escribir Ruta"
        placeholder="/Users/usuario/Proyectos/mi-proyecto"
        value={customPath}
        onChange={setCustomPath}
      />
    </Form>
  );
}
