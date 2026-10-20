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
  setProjectName,
  type FolttyProject,
} from "./core/projects";
import {
  getActiveTerminalName,
  TERMINAL_REGISTRY,
  type TerminalId,
} from "./core/terminal";

export default function Command() {
  const [projects, setProjects] = useState<FolttyProject[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadData = useCallback(() => {
    setIsLoading(true);
    setProjects(getProjects());
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
          description="Press ⌘+N to add a project from your Mac."
          actions={
            <ActionPanel>
              <Action.Push
                title="Add Project"
                icon={Icon.Plus}
                target={<AddProjectForm onProjectAdded={loadData} />}
              />
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
                  </ActionPanel.Section>

                  <ActionPanel.Section title="Project Management">
                    <Action.Push
                      title="Add Project"
                      icon={Icon.Plus}
                      shortcut={Keyboard.Shortcut.Common.New}
                      target={<AddProjectForm onProjectAdded={loadData} />}
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

function AddProjectForm({ onProjectAdded }: { onProjectAdded: () => void }) {
  const { pop } = useNavigation();
  const [folderPath, setFolderPath] = useState<string[]>([]);
  const [customPath, setCustomPath] = useState<string>("");

  const handleFolderChange = (paths: string[]) => {
    setFolderPath(paths);
    if (paths.length > 0 && paths[0]) {
      setCustomPath(paths[0]);
    }
  };

  const handleSubmit = async () => {
    const target =
      customPath.trim() || (folderPath.length > 0 ? folderPath[0] : "");

    if (!target) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Please choose a folder or enter a path",
      });
      return;
    }

    const res = addProject(target);
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
      message: path.basename(target),
    });

    onProjectAdded();
    pop();
  };

  return (
    <Form
      actions={
        <ActionPanel>
          <Action.SubmitForm
            title="Save Project"
            onSubmit={handleSubmit}
            shortcut={{ modifiers: [], key: "return" }}
          />
        </ActionPanel>
      }
    >
      <Form.Description text="Pick any folder from your Mac or enter a path." />

      <Form.FilePicker
        id="folderPicker"
        title="Choose Folder"
        canChooseDirectories={true}
        canChooseFiles={false}
        allowMultipleSelection={false}
        value={folderPath}
        onChange={handleFolderChange}
      />

      <Form.TextField
        id="customPath"
        title="Folder Path"
        placeholder="/Users/.../project or ~/Documents/code"
        value={customPath}
        onChange={setCustomPath}
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
    if (!res.success) {
      await showToast({
        style: Toast.Style.Failure,
        title: "Failed to Rename",
        message: res.message,
      });
      return;
    }

    await showToast({
      style: Toast.Style.Success,
      title: "Project Renamed",
      message: displayName.trim() || project.name,
    });

    onRenamed();
    pop();
  };

  return (
    <Form
      actions={
        <ActionPanel>
          <Action.SubmitForm
            title="Save Name"
            onSubmit={handleSubmit}
            shortcut={{ modifiers: [], key: "return" }}
          />
        </ActionPanel>
      }
    >
      <Form.Description
        text={`Set a friendly display name for ${project.path}`}
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
