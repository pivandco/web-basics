document.addEventListener("keydown", (event) => {
  if (event.code !== "Space") {
    return;
  }

  event.preventDefault();

  if (gameMode === GameMode.Playing) {
    setGameMode(GameMode.Paused);
  } else if (gameMode === GameMode.Paused) {
    setGameMode(GameMode.Playing);
  }
});

const GameMode = Object.freeze({
  Menu: "menu",
  Playing: "playing",
  Paused: "paused",
  GameOver: "gameover",
});

const ELEMENTS = Object.freeze({
  MENU: document.querySelector("#menu"),
  GAME: document.querySelector("#game"),
  FISHES: document.querySelector("#fishes"),
  PAUSE: document.querySelector("#pause"),
  GAME_OVER: document.querySelector("#game-over"),
  FINAL_SCORE: document.querySelector("#final-score"),
  SCOREBOARD: document.querySelector("#scoreboard"),
  SCOREBOARD_TBODY: document.querySelector("#scoreboard tbody"),
  SCOREBOARD_LOADING: document.querySelector("#scoreboard-loading"),
  SCOREBOARD_LOADING_FAILED: document.querySelector(
    "#scoreboard-loading-failed"
  ),
  NAME_TEXT_FIELD: document.querySelector("input#name"),
  PLAY_BUTTON: document.querySelector("button#play"),
  PLAYER_NAME: document.querySelector("#name-label"),
  SCORE: document.querySelector("#score"),
  TIME: document.querySelector("#time"),
  CLICK_SCORE: document.querySelector("#click-score"),
});

setInterval(() => {
  const clickScoreY = parseFloat(ELEMENTS.CLICK_SCORE.style.top);
  ELEMENTS.CLICK_SCORE.style.top = clickScoreY - 0.1 + "%";
  ELEMENTS.CLICK_SCORE.style.opacity -= 0.01;
}, 10);

let gameMode = GameMode.Menu;

const createGameState = () => ({
  time: 60,
  timer: setInterval(tickTime, 1000),
  fishInterval: setInterval(tickFishes, 10),
  score: 0,
});

function destroyGameState() {
  clearInterval(gameState.timer);
  clearInterval(gameState.fishInterval);
  Object.values(fishes).forEach((fish) => fish.kill());
}

let playerName;

const isTester = () => playerName.toLowerCase() === "tester";

let gameState;

function drawTimeAndScore() {
  if (gameState.time <= 10) {
    ELEMENTS.TIME.classList.add("red");
  }
  const minutes = Math.floor(gameState.time / 60)
    .toString()
    .padStart(2, "0");
  const seconds = (gameState.time % 60).toString().padStart(2, "0");
  ELEMENTS.TIME.textContent = `${minutes}:${seconds}`;

  ELEMENTS.SCORE.textContent = `Очки: ${gameState.score}`;
}

function tickFishes() {
  if (gameMode !== GameMode.Playing) {
    return;
  }
  Object.values(fishes).forEach((fish) => {
    fish.tick();
  });
}

function tickTime() {
  if (gameMode !== GameMode.Playing) {
    return;
  }

  if (!isTester()) {
    gameState.time--;
    if (gameState.time <= 0) {
      setGameMode(GameMode.GameOver);
    }
  }

  drawTimeAndScore();
  spawnFish();
}

const fishes = {};

const FISH_SPEED = 0.002;
const MAX_TARGET_CHANGE_TIME = 2000;
class Fish {
  #size;
  #targetChangesLeftUntilDeath = 4;
  #targetChangeTimeout;
  #target;
  #dying = false;
  #justSpawned;

  constructor(id) {
    this.id = id;
    this.#size = Math.round(Math.random() * 2 + 1);
    this.coords = this.#randomStartCoords();
    this.element = this.#createElement();
    this.#changeTarget();
    this.#justSpawned = true;
    this.#updateElement();
  }

  #randomStartCoords() {
    const coords = [Math.random() * 100, Math.random() * 100];
    if (Math.random() > 0.5) {
      coords[0] = Math.random() > 0.5 ? -10 : 110;
    } else {
      coords[1] = Math.random() > 0.5 ? -10 : 110;
    }
    return coords;
  }

  #changeTarget() {
    if (this.#dying) {
      return;
    }

    this.#justSpawned = false;

    this.#targetChangeTimeout = setTimeout(() => {
      this.#changeTarget();
    }, Math.random() * MAX_TARGET_CHANGE_TIME);

    if (gameMode !== GameMode.Playing) {
      return;
    }

    this.#target = [Math.random() * 100, Math.random() * 100];

    this.#targetChangesLeftUntilDeath--;
    if (this.#targetChangesLeftUntilDeath === 0) {
      this.#dying = true;
      this.#target = [120, 120];
    }
  }

  #createElement() {
    const element = document.createElement("img");

    element.id = this.id;
    element.src = "images/fish.png";
    element.classList.add("fish");
    element.classList.add(`fish-${this.#size}`);
    element.style.filter = `hue-rotate(${Math.random() * 360}deg)`;
    element.draggable = false;

    element.addEventListener("click", () => {
      handleFishCatch(this);
    });

    return element;
  }

  #updateElement() {
    if (gameMode !== GameMode.Playing) {
      return;
    }

    const [x, y] = this.coords;
    this.element.style.left = x + "%";
    this.element.style.top = y + "%";

    const [tx, ty] = this.#target;
    const angle = Math.atan2(ty - y, tx - x);
    this.element.style.transform = `rotate(${angle}rad)`;
  }

  tick() {
    if (gameMode !== GameMode.Playing) {
      return;
    }

    const [x, y] = this.coords;
    if ((x > 100 || y > 100) && !this.#justSpawned) {
      this.kill();
      return;
    }
    const [tx, ty] = this.#target;
    this.coords = [x + (tx - x) * FISH_SPEED, y + (ty - y) * FISH_SPEED];
    this.#updateElement();
  }

  get score() {
    return (4 - this.#size) * 10;
  }

  catch() {
    this.kill();
    ELEMENTS.CLICK_SCORE.style.display = "block";
    ELEMENTS.CLICK_SCORE.style.left = this.element.style.left;
    ELEMENTS.CLICK_SCORE.style.top = this.element.style.top;
    ELEMENTS.CLICK_SCORE.style.opacity = 1;
    ELEMENTS.CLICK_SCORE.textContent = `+${this.score}`;
  }

  kill() {
    clearTimeout(this.#targetChangeTimeout);
    this.element.remove();
  }
}

function spawnFish() {
  const id = Math.random().toString(36).slice(2, 9);
  const fish = new Fish(id);
  fishes[id] = fish;
  ELEMENTS.FISHES.appendChild(fish.element);
}

function handleFishCatch(fish) {
  fish.catch();
  gameState.score += fish.score;
  delete fishes[fish.id];
  drawTimeAndScore();
}

function setGameMode(newState) {
  const prevMode = gameMode;
  gameMode = newState;

  function show(el) {
    if (el.id === "scoreboard") {
      el.style.display = "";
    } else {
      el.style.display = "flex";
    }
  }
  function hide(el) {
    el.style.display = "none";
  }

  switch (gameMode) {
    case GameMode.Menu:
      show(ELEMENTS.MENU);
      hide(ELEMENTS.GAME);
      destroyGameState();
      break;
    case GameMode.Playing:
      show(ELEMENTS.GAME);
      hide(ELEMENTS.MENU);
      hide(ELEMENTS.PAUSE);

      ELEMENTS.PLAYER_NAME.textContent = playerName;
      if (isTester()) {
        ELEMENTS.PLAYER_NAME.classList.add("tester");
      }

      if (prevMode !== GameMode.Paused) {
        gameState = createGameState();
      }
      drawTimeAndScore();
      break;
    case GameMode.Paused:
      show(ELEMENTS.PAUSE);
      break;
    case GameMode.GameOver:
      show(ELEMENTS.GAME_OVER);
      show(ELEMENTS.SCOREBOARD_LOADING);
      hide(ELEMENTS.SCOREBOARD_LOADING_FAILED);
      hide(ELEMENTS.GAME);
      ELEMENTS.FINAL_SCORE.textContent = `Вы заработали ${gameState.score} очков`;
      destroyGameState();
      sendHighScore()
        .then(() => fillScoreboard())
        .then(() => {
          show(ELEMENTS.SCOREBOARD);
        })
        .catch((e) => {
          console.error("Failed to load scoreboard", e);
          show(ELEMENTS.SCOREBOARD_LOADING_FAILED);
        })
        .finally(() => {
          hide(ELEMENTS.SCOREBOARD_LOADING);
        });
      break;
  }
}

function handleNameFieldChange(event) {
  if (event.key === "Enter") {
    handlePlayButtonClick();
  }
  setTimeout(() => {
    ELEMENTS.PLAY_BUTTON.disabled = ELEMENTS.NAME_TEXT_FIELD.value.length === 0;
  });
}

function handlePlayButtonClick() {
  playerName = ELEMENTS.NAME_TEXT_FIELD.value;
  setGameMode(GameMode.Playing);
}

function handlePauseButtonClick() {
  setGameMode(GameMode.Paused);
}

function handleResumeButtonClick() {
  setGameMode(GameMode.Playing);
}

async function fillScoreboard() {
  const scores = await getScores();
  const table = ELEMENTS.SCOREBOARD_TBODY;
  table.innerHTML = "";
  scores.forEach((score) => {
    const row = document.createElement("tr");
    const ratingCell = document.createElement("td");
    const nameCell = document.createElement("td");
    const scoreCell = document.createElement("td");
    ratingCell.textContent = score.rating;
    nameCell.textContent = score.name;
    scoreCell.textContent = score.score;
    row.appendChild(ratingCell);
    row.appendChild(nameCell);
    row.appendChild(scoreCell);
    if (score.name === playerName) {
      row.classList.add("me");
    }
    table.appendChild(row);
  });
}

async function sendHighScore() {
  const response = await fetch("/api/high-scores", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: playerName,
      score: gameState.score,
    }),
  });
  if (!response.ok) {
    throw new Error("Failed to send high score");
  }
}

async function getScores() {
  const response = await fetch(
    "/api/high-scores?" + new URLSearchParams({ myName: playerName })
  );
  return await response.json();
}
