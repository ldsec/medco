<?php
$dsn = "pgsql:host=".getenv('I2B2_DB_HOST').";port=".getenv('I2B2_DB_PORT').";dbname=".getenv('I2B2_DB_NAME');
$options = [
  PDO::ATTR_EMULATE_PREPARES   => false, // turn off emulation mode for "real" prepared statements
  PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION, //turn on errors in the form of exceptions
];
try {
  $pdo = new PDO($dsn, getenv('I2B2_DB_USER'), getenv('I2B2_DB_PW'), $options);
} catch (Exception $e) {
  error_log($e->getMessage());
  exit('Error while connecting to the postgres database.'); //something a user can understand
}
?> 
