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
  getAllPresets,
  addCustomPreset,
  type FolttyProject,
  type FolttyPreset,
} from "./core/projects";

export default function Command() {
  const [projects, setProjects] = useState<FolttyProject[]>([]);
  const [presets, setPresets] = useState<FolttyPreset[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadData = useCallback(() => {
    setIsLoading(true);
    setProjects(getProjects());
    setPresets(getAllPresets());
    setIsLoading(false);
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

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

  const handleAddPreset = async (preset: FolttyPreset) => {
    const res = addProject(preset.path);
    if (res.success) {
      loadData();
      await showToast({
        style: Toast.Style.Success,
        title: "Preset Agregado",
        message: `${preset.name} -> ${preset.path}`,
      });
    } else {
      await showToast({
        style: Toast.Style.Failure,
        title: "No se pudo agregar",
        message: res.message,
      });
    }
  };

  const handleDelete = async (project: FolttyProject) => {
    const res = removeProject(project.path);
    if (res.success) {
      loadData();
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
          description="Agrega un proyecto con ⌘+N o usa uno de los accesos rápidos a continuación."
          actions={
            <ActionPanel>
              <Action.Push
                title="Agregar Proyecto Manual"
                icon={Icon.Plus}
                target={
                  <AddProjectForm onProjectAdded={loadData} presets={presets} />
                }
              />
              <ActionPanel.Section title="Presets Rápidos">
                {presets.slice(0, 4).map((p) => (
                  <Action
                    key={p.id}
                    title={`Agregar ${p.name}`}
                    icon={Icon.Star}
                    onAction={() => handleAddPreset(p)}
                  />
                ))}
              </ActionPanel.Section>
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

                <ActionPanel.Section title="Presets">
                  <ActionPanel.Submenu
                    title="Agregar Desde Preset…"
                    icon={Icon.Star}
                    shortcut={{ modifiers: ["cmd", "shift"], key: "a" }}
                  >
                    {presets.map((p) => (
                      <Action
                        key={p.id}
                        title={`Agregar ${p.name}`}
                        icon={p.isCustom ? Icon.Bookmark : Icon.Folder}
                        onAction={() => handleAddPreset(p)}
                      />
                    ))}
                  </ActionPanel.Submenu>

                  <Action.Push
                    title="Guardar Como Preset Custom…"
                    icon={Icon.Bookmark}
                    target={
                      <SavePresetForm
                        defaultPath={project.path}
                        onPresetSaved={loadData}
                      />
                    }
                  />
                </ActionPanel.Section>

                <ActionPanel.Section title="Gestión de Proyectos">
                  <Action.Push
                    title="Agregar Proyecto"
                    icon={Icon.Plus}
                    shortcut={{ modifiers: ["cmd"], key: "n" }}
                    target={
                      <AddProjectForm
                        onProjectAdded={loadData}
                        presets={presets}
                      />
                    }
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
                    onAction={loadData}
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

function AddProjectForm({
  onProjectAdded,
  presets,
}: {
  onProjectAdded: () => void;
  presets: FolttyPreset[];
}) {
  const { pop } = useNavigation();
  const [folderPath, setFolderPath] = useState<string[]>([]);
  const [customPath, setCustomPath] = useState<string>("");
  const [selectedPreset, setSelectedPreset] = useState<string>("");

  const handleSubmit = async () => {
    let target = "";

    if (selectedPreset) {
      const match = presets.find((p) => p.id === selectedPreset);
      if (match) target = match.path;
    }

    if (!target && folderPath.length > 0) {
      target = folderPath[0];
    }

    if (!target && customPath.trim()) {
      target = customPath.trim();
    }

    if (!target) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Por favor selecciona un preset, carpeta o escribe una ruta",
      });
      return;
    }

    const res = addProject(target);
    if (!res.success) {
      await showToast({
        style: Toast.Style.Failure,
        title: "No se pudo agregar",
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
      <Form.Description text="Puedes elegir un preset preconfigurado o seleccionar una carpeta del sistema." />

      <Form.Dropdown
        id="presetDropdown"
        title="O Usar Preset"
        value={selectedPreset}
        onChange={setSelectedPreset}
      >
        <Form.Dropdown.Item
          value=""
          title="-- Ninguno (Elegir Carpeta Abajo) --"
        />
        {presets.map((p) => (
          <Form.Dropdown.Item
            key={p.id}
            value={p.id}
            title={`${p.name} (${p.path.replace(homedir(), "~")})`}
          />
        ))}
      </Form.Dropdown>

      <Form.Separator />

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
        title="O Escribir Ruta / Preset"
        placeholder="docs, ~, ~/Documents o /ruta/completa"
        value={customPath}
        onChange={setCustomPath}
      />
    </Form>
  );
}

function SavePresetForm({
  defaultPath,
  onPresetSaved,
}: {
  defaultPath: string;
  onPresetSaved: () => void;
}) {
  const { pop } = useNavigation();
  const [presetId, setPresetId] = useState<string>(
    path
      .basename(defaultPath)
      .toLowerCase()
      .replace(/[^a-z0-9_-]/g, ""),
  );
  const [targetPath, setTargetPath] = useState<string>(defaultPath);

  const handleSubmit = async () => {
    if (!presetId.trim() || !targetPath.trim()) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Por favor completa el ID y la ruta del preset",
      });
      return;
    }

    const res = addCustomPreset(presetId, targetPath);
    if (!res.success) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Error al guardar preset",
        message: res.message,
      });
      return;
    }

    await showToast({
      style: Toast.Style.Success,
      title: "Preset Guardado",
      message: `${presetId} -> ${targetPath}`,
    });

    onPresetSaved();
    pop();
  };

  return (
    <Form
      actions={
        <ActionPanel>
          <Action.SubmitForm title="Guardar Preset" onSubmit={handleSubmit} />
        </ActionPanel>
      }
    >
      <Form.Description text="Crea un alias para poder abrir o agregar esta ruta rápidamente (ej: 'work', 'blog')." />
      <Form.TextField
        id="presetId"
        title="Identificador (Alias)"
        placeholder="ej: work, proyectos, blog"
        value={presetId}
        onChange={setPresetId}
      />
      <Form.TextField
        id="presetPath"
        title="Ruta del Preset"
        value={targetPath}
        onChange={setTargetPath}
      />
    </Form>
  );
}
