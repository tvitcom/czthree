-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Dec 09, 2021 at 09:43 AM
-- Server version: 8.0.27
-- PHP Version: 7.3.31-1~deb10u1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `czthree`
--

-- --------------------------------------------------------

--
-- Table structure for table `todo`
--

CREATE TABLE `todo` (
  `todo_id` bigint UNSIGNED NOT NULL,
  `author_id` bigint NOT NULL DEFAULT '0',
  `perfomer_id` bigint NOT NULL DEFAULT '0',
  `name` varchar(128) NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0' COMMENT 'canceled:0, todo:1, done:2'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `todo`
--

INSERT INTO `todo` (`todo_id`, `author_id`, `perfomer_id`, `name`, `status`) VALUES
(1, 1, 1, 'ЗАпланировать задачи сотрудникам', 1),
(2, 1, 2, 'Высотные работы№1 по плану 1', 1),
(3, 1, 2, 'Высотные работы№2 по плану 1', 1),
(4, 1, 3, 'Подготовка высотных работ по плану 1', 1);

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `user_id` bigint UNSIGNED NOT NULL,
  `name` varchar(128) NOT NULL,
  `email` varchar(128) NOT NULL,
  `passhash` varchar(128) NOT NULL,
  `role` tinyint NOT NULL DEFAULT '0' COMMENT 'admin:100, user:10, guest:0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`user_id`, `name`, `email`, `passhash`, `role`) VALUES
(1, 'admin', 'admin@gmail.com', '$2a$12$RNw.W2JPNHRmlBbxKw5aiuDUW1OdojjOKb/68BXNsThqtc/jM2vf.', 100),
(2, 'demo', 'demo@mail.ru', '$2a$12$FnKg9DqpDeDzpThkouc1b.SbZy4qAqgKBGohswrWBlzr7jDDrKYOW', 10),
(3, 'Vaso', 'vaso@mail.ru', '$2a$12$dcWMn4UGnRcEhVsDeYpD1OWkc5fQpogXqUTyKwXZzKjPNAjbYA4xC', 10);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `todo`
--
ALTER TABLE `todo`
  ADD PRIMARY KEY (`todo_id`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`user_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `todo`
--
ALTER TABLE `todo`
  MODIFY `todo_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `user`
--
ALTER TABLE `user`
  MODIFY `user_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
