<?php
$content = $_POST['content'];
$come = base64_decode($content);
$process = json_encode($come,true);
$out = json_decode($process,true);
$file = fopen("ban.json", "w");
fwrite($file, $out);
fclose($file);
?>