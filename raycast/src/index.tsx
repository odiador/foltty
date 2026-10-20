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
  Keyboard,
} from "@raycast/api";
import { homedir } from "os";
import path from "path";
import { useState, useEffect, useCallback } from "react";
import {
  getProjects,
  addProject,
  removeProject,
  openProjectInTerminal,
  getStorageFilePath,
  getAllPresets,
  addCustomPreset,
  setProjectName,
  type FolttyProject,
  type FolttyPreset,
} from "./core/projects";
import {
  getActiveTerminalName,
  TERMINAL_REGISTRY,
  type TerminalId,
} from "./core/terminal";

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

  const activeTerminalName = getActiveTerminalName();

  const handleOpenTerminal = async (
    projectPath: string,
    targetTerminal?: TerminalId,
    ghosttyMode?: "tab" | "window",
  ) => {
    try {
      await closeMainWindow();
      await openProjectInTerminal(projectPath, targetTerminal, ghosttyMode);
    } catch (error) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Error Launching Terminal",
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
        title: "Preset Added",
        message: `${preset.name} -> ${preset.path}`,
      });
    } else {
      await showToast({
        style: Toast.Style.Failure,
        title: "Failed to Add",
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
        title: "Project Removed",
        message: project.name,
      });
    } else {
      await showToast({
        style: Toast.Style.Failure,
        title: "Failed to Remove",
        message: res.message,
      });
    }
  };

  const storageFile = getStorageFilePath();

  return (
    <List
      isLoading={isLoading}
      searchBarPlaceholder="Filter projects by name or path..."
    >
      {projects.length === 0 && !isLoading ? (
        <List.EmptyView
          icon="icon.png"
          title="No Projects Registered in Foltty"
          description="Press ⌘+N to add a project or pick a quick preset below."
          actions={
            <ActionPanel>
              <Action.Push
                title="Add Project"
                icon={Icon.Plus}
                target={
                  <AddProjectForm onProjectAdded={loadData} presets={presets} />
                }
              />
              <ActionPanel.Section title="Quick Presets">
                {presets.slice(0, 4).map((p) => (
                  <Action
                    key={p.id}
                    title={`Add ${p.name}`}
                    icon={Icon.Star}
                    onAction={() => handleAddPreset(p)}
                  />
                ))}
              </ActionPanel.Section>
              <Action.OpenWith
                title="Open Projects File"
                path={storageFile}
                shortcut={{ modifiers: ["cmd", "shift"], key: "p" }}
              />
            </ActionPanel>
          }
        />
      ) : (
        projects.map((project) => {
          const accessories: List.Item.Accessory[] = [];
          if (!project.exists) {
            accessories.push({
              tag: { value: "Missing Path", color: Color.Red },
              tooltip: "Folder does not exist on disk",
            });
          }
          if (project.customName) {
            accessories.push({
              tag: { value: "Custom Name", color: Color.Purple },
            });
          }

          return (
            <List.Item
              key={project.id}
              icon={
                project.exists
                  ? "icon.png"
                  : { source: Icon.ExclamationMark, tintColor: Color.Red }
              }
              title={project.name}
              subtitle={project.path.replace(homedir(), "~")}
              accessories={accessories}
              actions={
                <ActionPanel>
                  <ActionPanel.Section title="Terminal">
                    <Action
                      title={`Open in ${activeTerminalName}`}
                      icon={Icon.Terminal}
                      onAction={() => handleOpenTerminal(project.path)}
                    />
                    {TERMINAL_REGISTRY.ghostty.isAvailable() && (
                      <>
                        <Action
                          title="Open in Ghostty (New Tab)"
                          icon={Icon.PlusCircle}
                          shortcut={{ modifiers: ["cmd"], key: "return" }}
                          onAction={() =>
                            handleOpenTerminal(project.path, "ghostty", "tab")
                          }
                        />
                        <Action
                          title="Open in Ghostty (New Window)"
                          icon={Icon.Window}
                          shortcut={{
                            modifiers: ["cmd", "shift"],
                            key: "return",
                          }}
                          onAction={() =>
                            handleOpenTerminal(
                              project.path,
                              "ghostty",
                              "window",
                            )
                          }
                        />
                      </>
                    )}
                    <ActionPanel.Submenu
                      title="Open in Specific Terminal…"
                      icon={Icon.Terminal}
                      shortcut={{ modifiers: ["cmd"], key: "t" }}
                    >
                      <Action
                        title="Ghostty (New Tab)"
                        icon={Icon.PlusCircle}
                        onAction={() =>
                          handleOpenTerminal(project.path, "ghostty", "tab")
                        }
                      />
                      <Action
                        title="Ghostty (New Window)"
                        icon={Icon.Window}
                        onAction={() =>
                          handleOpenTerminal(project.path, "ghostty", "window")
                        }
                      />
                      {(
                        Object.keys(
                          TERMINAL_REGISTRY,
                        ) as (keyof typeof TERMINAL_REGISTRY)[]
                      )
                        .filter((termKey) => termKey !== "ghostty")
                        .map((termKey) => {
                          const term = TERMINAL_REGISTRY[termKey];
                          const available = term.isAvailable();
                          return (
                            <Action
                              key={term.id}
                              title={`${term.name}${available ? "" : " (Not Installed)"}`}
                              icon={
                                available ? Icon.Terminal : Icon.XMarkCircle
                              }
                              onAction={() =>
                                handleOpenTerminal(project.path, term.id)
                              }
                            />
                          );
                        })}
                    </ActionPanel.Submenu>
                  </ActionPanel.Section>

                  <ActionPanel.Section title="File Actions">
                    <Action.ShowInFinder path={project.path} />
                    <Action.OpenWith
                      path={project.path}
                      shortcut={Keyboard.Shortcut.Common.OpenWith}
                    />
                    <Action.CopyToClipboard
                      title="Copy Path"
                      content={project.path}
                      shortcut={{ modifiers: ["cmd"], key: "c" }}
                    />
                  </ActionPanel.Section>

                  <ActionPanel.Section title="Project Details">
                    <Action.Push
                      title="Rename Display Name…"
                      icon={Icon.Pencil}
                      shortcut={Keyboard.Shortcut.Common.Refresh}
                      target={
                        <RenameProjectForm
                          project={project}
                          onRenamed={loadData}
                        />
                      }
                    />
                    <Action.Push
                      title="Save as Custom Preset…"
                      icon={Icon.Bookmark}
                      target={
                        <SavePresetForm
                          defaultPath={project.path}
                          onPresetSaved={loadData}
                        />
                      }
                    />
                  </ActionPanel.Section>

                  <ActionPanel.Section title="Presets">
                    <ActionPanel.Submenu
                      title="Add from Preset…"
                      icon={Icon.Star}
                      shortcut={{ modifiers: ["cmd", "shift"], key: "a" }}
                    >
                      {presets.map((p) => (
                        <Action
                          key={p.id}
                          title={`Add ${p.name}`}
                          icon={p.isCustom ? Icon.Bookmark : Icon.Folder}
                          onAction={() => handleAddPreset(p)}
                        />
                      ))}
                    </ActionPanel.Submenu>
                  </ActionPanel.Section>

                  <ActionPanel.Section title="Project Management">
                    <Action.Push
                      title="Add Project"
                      icon={Icon.Plus}
                      shortcut={Keyboard.Shortcut.Common.New}
                      target={
                        <AddProjectForm
                          onProjectAdded={loadData}
                          presets={presets}
                        />
                      }
                    />
                    <Action
                      title="Remove from List"
                      icon={Icon.Trash}
                      style={Action.Style.Destructive}
                      shortcut={{ modifiers: ["cmd"], key: "backspace" }}
                      onAction={() => handleDelete(project)}
                    />
                    <Action
                      title="Reload List"
                      icon={Icon.ArrowClockwise}
                      onAction={loadData}
                    />
                  </ActionPanel.Section>
                </ActionPanel>
              }
            />
          );
        })
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
  const [customName, setCustomName] = useState<string>("");
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
        title: "Please choose a preset, pick a folder, or enter a path",
      });
      return;
    }

    const res = addProject(target, customName.trim());
    if (!res.success) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Failed to Add Project",
        message: res.message,
      });
      return;
    }

    await showToast({
      style: Toast.Style.Success,
      title: "Project Added",
      message: customName || path.basename(target),
    });

    onProjectAdded();
    pop();
  };

  return (
    <Form
      actions={
        <ActionPanel>
          <Action.SubmitForm title="Save Project" onSubmit={handleSubmit} />
        </ActionPanel>
      }
    >
      <Form.Description text="Choose a predefined preset or pick any folder from your Mac." />

      <Form.Dropdown
        id="presetDropdown"
        title="Use Preset"
        value={selectedPreset}
        onChange={setSelectedPreset}
      >
        <Form.Dropdown.Item value="" title="-- None (Pick Folder Below) --" />
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
        title="Choose Folder"
        canChooseDirectories={true}
        canChooseFiles={false}
        allowMultipleSelection={false}
        value={folderPath}
        onChange={setFolderPath}
      />

      <Form.TextField
        id="customPath"
        title="Or Enter Path / Preset"
        placeholder="docs, ~, ~/Documents or /full/path"
        value={customPath}
        onChange={setCustomPath}
      />

      <Form.TextField
        id="customName"
        title="Display Name (Optional)"
        placeholder="e.g. Backend API, Personal Blog"
        value={customName}
        onChange={setCustomName}
      />
    </Form>
  );
}

function RenameProjectForm({
  project,
  onRenamed,
}: {
  project: FolttyProject;
  onRenamed: () => void;
}) {
  const { pop } = useNavigation();
  const [displayName, setDisplayName] = useState(
    project.customName || project.name,
  );

  const handleSubmit = async () => {
    const res = setProjectName(project.path, displayName.trim());
    if (res.success) {
      await showToast({
        style: Toast.Style.Success,
        title: "Project Renamed",
        message: displayName || path.basename(project.path),
      });
      onRenamed();
      pop();
    } else {
      await showToast({
        style: Toast.Style.Failure,
        title: "Error Renaming",
        message: res.message,
      });
    }
  };

  return (
    <Form
      actions={
        <ActionPanel>
          <Action.SubmitForm
            title="Save Display Name"
            onSubmit={handleSubmit}
          />
        </ActionPanel>
      }
    >
      <Form.Description
        text={`Assign a custom friendly display name for ${project.path}`}
      />
      <Form.TextField
        id="displayName"
        title="Display Name"
        value={displayName}
        onChange={setDisplayName}
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
        title: "Please provide both an alias ID and a target path",
      });
      return;
    }

    const res = addCustomPreset(presetId, targetPath);
    if (!res.success) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Failed to Save Preset",
        message: res.message,
      });
      return;
    }

    await showToast({
      style: Toast.Style.Success,
      title: "Preset Saved",
      message: `${presetId} -> ${targetPath}`,
    });

    onPresetSaved();
    pop();
  };

  return (
    <Form
      actions={
        <ActionPanel>
          <Action.SubmitForm title="Save Preset" onSubmit={handleSubmit} />
        </ActionPanel>
      }
    >
      <Form.Description text="Create a quick alias for this directory (e.g. 'work', 'blog', 'sandbox')." />
      <Form.TextField
        id="presetId"
        title="Alias / Preset ID"
        placeholder="e.g. work, projects, blog"
        value={presetId}
        onChange={setPresetId}
      />
      <Form.TextField
        id="presetPath"
        title="Preset Path"
        value={targetPath}
        onChange={setTargetPath}
      />
    </Form>
  );
}
