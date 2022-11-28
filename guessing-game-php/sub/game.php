<?php

const MIN_ANSWER = 1;
const MAX_ANSWER = 100;

session_start();

if ($_SERVER['REQUEST_METHOD'] == 'GET') {
    $_SESSION = [
        'actual_number' => rand(MIN_ANSWER, MAX_ANSWER),
        'attempts' => ceil(log(MAX_ANSWER - MIN_ANSWER + 1, 2)),
        'attempts_log' => [],
    ];
}

$actual_number = $_SESSION['actual_number'];
$message = '';
$game_over = false;

if ($_SERVER['REQUEST_METHOD'] == 'POST') {
    $answer = $_POST['answer'];
    if (!is_numeric($answer)) {
        $message = 'Введите число';
    } elseif ($answer < MIN_ANSWER || $answer > MAX_ANSWER) {
        $message = 'Число должно быть в диапазоне от ' . MIN_ANSWER . ' до ' . MAX_ANSWER;
    } else {
        $_SESSION['attempts']--;
        $_SESSION['attempts_log'][] = $answer;
        if ($_SESSION['attempts'] == 0) {
            $message = "К сожалению, вы проиграли! Мое число было {$actual_number}. Удачи в следующий раз!";
            $game_over = true;
        } elseif ($answer == $actual_number) {
            $message = 'Вы угадали! 🎉';
            $game_over = true;
        } elseif ($answer < $actual_number) {
            $message = 'Не угадали, мое число больше';
        } elseif ($answer > $actual_number) {
            $message = 'Не угадали, мое число меньше';
        }
    }
}

$formatted_attempts_log = '';
foreach ($_SESSION['attempts_log'] as $attempt) {
    if ($attempt < $actual_number) {
        $formatted_attempts_log .= "<li>{$attempt} (нужно больше)</li>";
    } elseif ($attempt > $actual_number) {
        $formatted_attempts_log .= "<li>{$attempt} (нужно меньше)</li>";
    } else {
        $formatted_attempts_log .= "<li>{$attempt} (вы угадали)</li>";
    }
}
