<?php require 'game.php'; ?>

<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Угадай число</title>
    <link rel="stylesheet" href="/game.css"/>
</head>

<body>

<div class="game">
    <h1>Игра "Угадай число"</h1>
    Компьютер загадал число в диапазоне [<?= MIN_ANSWER ?>; <?= MAX_ANSWER ?>]
    <p>Попыток: <span><?= $_SESSION['attempts'] ?></span></p>
    <hr>
    <?php if ($formatted_attempts_log): ?>
        <h2>Ваши предыдущие попытки</h2>
        <ol id="attempts"><?= $formatted_attempts_log ?></ol>
        <hr>
    <?php endif ?>
    <p><?= $message ?></p>

    <form method="post" class="controls">
        <input type="number" name="answer" placeholder="Число, пожалуйста" autofocus <?= $game_over ? 'disabled' : '' ?> />
        <div class="buttons">
            <button class="primary" type="submit" <?= $game_over ? 'disabled' : '' ?>>Проверить</button>
            <a href="<?= $_SERVER['REQUEST_URI'] ?>" class="primary">Заново</a>
        </div>
    </form>
</div>
</body>

</html>
