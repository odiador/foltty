import { spawn, spawnSync } from "child_process";
import readline from "readline";
import {
  getProjects,
  addProject,
  removeProject,
  openInGhostty,
  getStorageFilePath,
} from "./core/projects";

function printHelp() {
  console.log(`
👻 Foltty - Project Manager para Ghostty & Raycast

USO:
  foltty                     Abre la TUI interactiva
  foltty add [ruta]          Agrega la carpeta actual o una ruta a la lista
  foltty rm <ruta|nombre>    Elimina un proyecto de la lista
  foltty list                Muestra todos los proyectos registrados
  foltty open <nombre>       Abre directamente un proyecto por nombre o búsqueda
  foltty file                Muestra la ruta del archivo de proyectos
  foltty --help              Muestra esta ayuda

ATAJOS EN LA TUI:
  [ENTER]   Abrir proyecto en Ghostty
  [CTRL-X]  Eliminar de la lista
  [ESC]     Salir
`);
}

function hasFzf(): boolean {
  try {
    const res = spawnSync("which", ["fzf"], { stdio: "ignore" });
    return res.status === 0;
  } catch {
    return false;
  }
}

/**
 * TUI usando FZF con theming moderno y preview
 */
function runFzfTui() {
  const projects = getProjects();
  if (projects.length === 0) {
    console.log("👻 No hay proyectos registrados en Foltty.");
    console.log("👉 Agrega uno con: foltty add");
    return;
  }

  const storageFile = getStorageFilePath();
  const inputList = projects.map((p) => p.path).join("\n");

  const fzfArgs = [
    "--prompt=👻 Foltty > ",
    "--height=70%",
    "--layout=reverse",
    "--border=rounded",
    "--border-label= 📂 FOLTTY - GHOSTTY SELECTOR ",
    "--color=border:#8b5cf6,label:#ec4899,prompt:#38bdf8,pointer:#f43f5e",
    "--header=[ENTER] Abrir en Ghostty  |  [CTRL-X] Eliminar",
    "--header-first",
    `--bind=ctrl-x:execute(node -e "require('${__dirname}/core/projects').removeProject('{}')")+reload(cat ${storageFile})`,
    "--preview=echo '📁 Ruta: {}'; echo ''; if [ -d '{}' ]; then ls -la --color=always '{}' | head -n 15; else echo '❌ Ruta no encontrada'; fi",
    "--preview-window=right:50%:border-left",
  ];

  const fzf = spawn("fzf", fzfArgs, {
    stdio: ["pipe", "pipe", "inherit"],
  });

  fzf.stdin.write(inputList);
  fzf.stdin.end();

  let selected = "";
  fzf.stdout.on("data", (data) => {
    selected += data.toString();
  });

  fzf.on("close", async (code) => {
    const chosen = selected.trim();
    if (code === 0 && chosen) {
      try {
        await openInGhostty(chosen);
        console.log(`🚀 Abriendo en Ghostty: ${chosen}`);
      } catch (err) {
        console.error(
          "Error al abrir Ghostty:",
          err instanceof Error ? err.message : err,
        );
      }
    }
  });
}

/**
 * TUI nativa interactiva en Node (fallback si no hay fzf)
 */
function runNativeTui() {
  const projects = getProjects();
  if (projects.length === 0) {
    console.log("👻 No hay proyectos registrados en Foltty.");
    console.log("👉 Agrega uno ejecutando: foltty add");
    return;
  }

  let selectedIndex = 0;
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
  });

  readline.emitKeypressEvents(process.stdin);
  if (process.stdin.isTTY) {
    process.stdin.setRawMode(true);
  }

  function render() {
    console.clear();
    console.log("\x1b[35m👻 FOLTTY - Project Manager\x1b[0m");
    console.log(
      "\x1b[90mUsa ↑/↓ para navegar, [ENTER] para abrir, [d] para eliminar, [q] para salir\x1b[0m\n",
    );

    projects.forEach((proj, idx) => {
      const isSelected = idx === selectedIndex;
      const pointer = isSelected ? "\x1b[36m❯ \x1b[0m" : "  ";
      const name = isSelected
        ? `\x1b[1m\x1b[32m${proj.name}\x1b[0m`
        : proj.name;
      const status = proj.exists ? "" : " \x1b[31m[Inaccesible]\x1b[0m";
      console.log(`${pointer}${name} \x1b[90m(${proj.path})\x1b[0m${status}`);
    });
  }

  render();

  process.stdin.on("keypress", async (_, key) => {
    if (key.ctrl && key.name === "c") {
      process.exit(0);
    }

    if (key.name === "up") {
      selectedIndex = (selectedIndex - 1 + projects.length) % projects.length;
      render();
    } else if (key.name === "down") {
      selectedIndex = (selectedIndex + 1) % projects.length;
      render();
    } else if (key.name === "return") {
      const target = projects[selectedIndex];
      if (process.stdin.isTTY) {
        process.stdin.setRawMode(false);
      }
      rl.close();
      try {
        await openInGhostty(target.path);
        console.log(`\n🚀 Abriendo en Ghostty: ${target.path}`);
      } catch (err) {
        console.error(
          "\n❌ Error al abrir Ghostty:",
          err instanceof Error ? err.message : err,
        );
      }
      process.exit(0);
    } else if (key.name === "d") {
      const target = projects[selectedIndex];
      removeProject(target.path);
      projects.splice(selectedIndex, 1);
      if (selectedIndex >= projects.length) {
        selectedIndex = Math.max(0, projects.length - 1);
      }
      if (projects.length === 0) {
        if (process.stdin.isTTY) {
          process.stdin.setRawMode(false);
        }
        rl.close();
        console.clear();
        console.log("No quedan proyectos registrados.");
        process.exit(0);
      }
      render();
    } else if (key.name === "q" || key.name === "escape") {
      if (process.stdin.isTTY) {
        process.stdin.setRawMode(false);
      }
      rl.close();
      console.clear();
      process.exit(0);
    }
  });
}

export async function runCli(args: string[]): Promise<void> {
  const command = args[0];

  if (!command || command === "tui") {
    if (hasFzf()) {
      runFzfTui();
    } else {
      runNativeTui();
    }
    return;
  }

  if (command === "add") {
    const target = args[1] || process.cwd();
    const res = addProject(target);
    if (res.success) {
      console.log(`✅ ${res.message}`);
    } else {
      console.error(`❌ ${res.message}`);
      process.exit(1);
    }
    return;
  }

  if (command === "rm" || command === "remove" || command === "delete") {
    const target = args[1];
    if (!target) {
      console.error("❌ Especifica la ruta o nombre del proyecto a eliminar.");
      process.exit(1);
    }
    const res = removeProject(target);
    if (res.success) {
      console.log(`🗑️ ${res.message}`);
    } else {
      console.error(`❌ ${res.message}`);
      process.exit(1);
    }
    return;
  }

  if (command === "list" || command === "ls") {
    const projects = getProjects();
    if (projects.length === 0) {
      console.log("👻 No hay proyectos registrados.");
      return;
    }
    console.log(`👻 Proyectos registrados en Foltty (${projects.length}):\n`);
    projects.forEach((p) => {
      const mark = p.exists ? "📂" : "⚠️ ";
      console.log(
        `  ${mark} \x1b[1m${p.name}\x1b[0m -> \x1b[90m${p.path}\x1b[0m`,
      );
    });
    return;
  }

  if (command === "open") {
    const query = args[1];
    if (!query) {
      console.error("❌ Especifica el nombre o ruta del proyecto a abrir.");
      process.exit(1);
    }
    const projects = getProjects();
    const match = projects.find(
      (p) =>
        p.name.toLowerCase() === query.toLowerCase() ||
        p.path.toLowerCase().includes(query.toLowerCase()),
    );

    if (!match) {
      console.error(
        `❌ No se encontró ningún proyecto con el criterio: ${query}`,
      );
      process.exit(1);
    }

    try {
      await openInGhostty(match.path);
      console.log(`🚀 Abriendo en Ghostty: ${match.path}`);
    } catch (err) {
      console.error(
        `❌ Error al abrir Ghostty:`,
        err instanceof Error ? err.message : err,
      );
      process.exit(1);
    }
    return;
  }

  if (command === "file") {
    console.log(getStorageFilePath());
    return;
  }

  if (command === "--help" || command === "-h" || command === "help") {
    printHelp();
    return;
  }

  console.error(`Comando desconocido: "${command}"\n`);
  printHelp();
  process.exit(1);
}

runCli(process.argv.slice(2));
