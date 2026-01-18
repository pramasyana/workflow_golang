-- phpMyAdmin SQL Dump
-- version 4.8.5
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Jan 18, 2026 at 10:00 AM
-- Server version: 10.1.38-MariaDB
-- PHP Version: 5.6.40

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `workflow_approval`
--

-- --------------------------------------------------------

--
-- Table structure for table `actors`
--

CREATE TABLE `actors` (
  `id` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `code` varchar(50) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `actors`
--

INSERT INTO `actors` (`id`, `name`, `code`, `created_at`, `updated_at`) VALUES
('50d3c341-469b-4ea4-bf57-bf7fd5894674', 'CEO', 'ceo', '2026-01-18 08:51:15', '2026-01-18 08:51:15'),
('9bd37fb3-3578-4353-bac0-e1fb83f0a7f6', 'Manager', 'manager', '2026-01-18 08:50:53', '2026-01-18 08:50:53'),
('f2777759-f96b-4b4c-82b0-e9e6e65dea39', 'Director', 'director', '2026-01-18 08:51:06', '2026-01-18 08:51:06');

-- --------------------------------------------------------

--
-- Table structure for table `approval_history`
--

CREATE TABLE `approval_history` (
  `id` varchar(36) NOT NULL,
  `request_id` varchar(36) NOT NULL,
  `workflow_id` varchar(36) NOT NULL,
  `step_level` int(11) NOT NULL,
  `actor_id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `action` varchar(20) NOT NULL,
  `comment` text,
  `created_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `approval_history`
--

INSERT INTO `approval_history` (`id`, `request_id`, `workflow_id`, `step_level`, `actor_id`, `user_id`, `action`, `comment`, `created_at`) VALUES
('09a4ed79-baff-41d3-80c2-f47e98797e36', '4c126105-06fd-4bc8-8d86-d6e235da820e', '1e876554-e267-498d-920e-f09b93d7ab32', 3, '50d3c341-469b-4ea4-bf57-bf7fd5894674', 'b0065105-7069-4764-8b68-9321ea313c1b', 'APPROVE', '', '2026-01-18 08:57:53'),
('4eef9d6e-13d8-4ffa-8ead-5fc13d3e7ed6', 'c1602ff4-61a5-4997-af0c-e0669303af0c', '1e876554-e267-498d-920e-f09b93d7ab32', 2, 'f2777759-f96b-4b4c-82b0-e9e6e65dea39', '9d963246-c42b-4e72-b45a-b1b73cef9330', 'REJECT', 'Insufficient documentation', '2026-01-18 08:59:22'),
('53ee3ae1-1608-4612-bb65-739cc1285235', 'c1602ff4-61a5-4997-af0c-e0669303af0c', '1e876554-e267-498d-920e-f09b93d7ab32', 1, '9bd37fb3-3578-4353-bac0-e1fb83f0a7f6', '008905a7-732e-43ec-a6e9-ef40424db70d', 'APPROVE', '', '2026-01-18 08:58:49'),
('cc33ed5a-ff07-4d43-80fa-adf48c9d059f', '4c126105-06fd-4bc8-8d86-d6e235da820e', '1e876554-e267-498d-920e-f09b93d7ab32', 2, 'f2777759-f96b-4b4c-82b0-e9e6e65dea39', '9d963246-c42b-4e72-b45a-b1b73cef9330', 'APPROVE', '', '2026-01-18 08:57:29'),
('f2eac32f-cced-4805-bfba-b0835d80dbfa', '4c126105-06fd-4bc8-8d86-d6e235da820e', '1e876554-e267-498d-920e-f09b93d7ab32', 1, '9bd37fb3-3578-4353-bac0-e1fb83f0a7f6', '008905a7-732e-43ec-a6e9-ef40424db70d', 'APPROVE', '', '2026-01-18 08:57:00');

-- --------------------------------------------------------

--
-- Table structure for table `requests`
--

CREATE TABLE `requests` (
  `id` varchar(36) NOT NULL,
  `workflow_id` varchar(36) NOT NULL,
  `requester_id` varchar(36) NOT NULL,
  `current_step` int(11) DEFAULT '1',
  `status` varchar(20) DEFAULT 'PENDING',
  `amount` decimal(15,2) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `description` text,
  `version` int(11) DEFAULT '1',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `requests`
--

INSERT INTO `requests` (`id`, `workflow_id`, `requester_id`, `current_step`, `status`, `amount`, `title`, `description`, `version`, `created_at`, `updated_at`) VALUES
('4c126105-06fd-4bc8-8d86-d6e235da820e', '1e876554-e267-498d-920e-f09b93d7ab32', 'b693fdee-f2bb-11f0-8cc1-7a447c8d071a', 4, 'APPROVED', '5000000.00', 'Office Supplies Purchase', 'Monthly office supplies for Q1', 4, '2026-01-18 08:55:56', '2026-01-18 08:57:53'),
('c1602ff4-61a5-4997-af0c-e0669303af0c', '1e876554-e267-498d-920e-f09b93d7ab32', 'b0065105-7069-4764-8b68-9321ea313c1b', 2, 'REJECTED', '5000000.00', 'Office Supplies Purchase', 'Monthly office supplies for Q1', 3, '2026-01-18 08:58:10', '2026-01-18 08:59:22');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `is_admin` tinyint(1) DEFAULT '0',
  `actor_id` varchar(36) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `email`, `password`, `name`, `is_admin`, `actor_id`, `created_at`, `updated_at`) VALUES
('008905a7-732e-43ec-a6e9-ef40424db70d', 'manager@gmail.com', '$2a$10$ZpHtOVq2QtL77JavbOD2T.wXp.6k/4ltRwhaT1M5m.841u0uT7bzi', 'Manager', 0, '9bd37fb3-3578-4353-bac0-e1fb83f0a7f6', '2026-01-18 08:51:54', '2026-01-18 08:51:54'),
('9d963246-c42b-4e72-b45a-b1b73cef9330', 'director@gmail.com', '$2a$10$sKPAbfAOzstPxE5DHSyUYOj5RazWsgKrOzaerM.KFGNBx.eAsBaQK', 'Director', 0, 'f2777759-f96b-4b4c-82b0-e9e6e65dea39', '2026-01-18 08:52:23', '2026-01-18 08:52:23'),
('b0065105-7069-4764-8b68-9321ea313c1b', 'ceo@gmail.com', '$2a$10$Ufkru0IT35EUXC1l3H6UAewr6wVCcRcpb5DUlPrGvNXRZJZvGL8oW', 'CEO', 0, '50d3c341-469b-4ea4-bf57-bf7fd5894674', '2026-01-18 08:52:45', '2026-01-18 08:52:45'),
('b693fdee-f2bb-11f0-8cc1-7a447c8d071a', 'administrator@gmail.com', '$2a$10$qToIHURdilgF4kZCTnvR5.8AVRuljFLW1vHoRpAWjDIQjpPZyZ.ie', 'Administrator', 1, NULL, '2026-01-18 15:50:06', '2026-01-18 15:50:06');

-- --------------------------------------------------------

--
-- Table structure for table `workflows`
--

CREATE TABLE `workflows` (
  `id` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `workflows`
--

INSERT INTO `workflows` (`id`, `name`, `created_at`, `updated_at`) VALUES
('1e876554-e267-498d-920e-f09b93d7ab32', 'Purchase Approval', '2026-01-18 08:53:14', '2026-01-18 08:53:14');

-- --------------------------------------------------------

--
-- Table structure for table `workflow_steps`
--

CREATE TABLE `workflow_steps` (
  `id` varchar(36) NOT NULL,
  `workflow_id` varchar(36) NOT NULL,
  `level` int(11) NOT NULL,
  `actor_id` varchar(36) NOT NULL,
  `conditions` text,
  `description` varchar(500) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `workflow_steps`
--

INSERT INTO `workflow_steps` (`id`, `workflow_id`, `level`, `actor_id`, `conditions`, `description`, `created_at`, `updated_at`) VALUES
('167d48ad-fe92-4bfe-a97d-cc583dc609e8', '1e876554-e267-498d-920e-f09b93d7ab32', 1, '9bd37fb3-3578-4353-bac0-e1fb83f0a7f6', '{\"min_amount\":1000000}', '', '2026-01-18 08:55:11', '2026-01-18 08:55:11'),
('323b83fb-69bc-444e-8f92-05bb7fd09a2f', '1e876554-e267-498d-920e-f09b93d7ab32', 2, 'f2777759-f96b-4b4c-82b0-e9e6e65dea39', '{\"min_amount\":1000000}', '', '2026-01-18 08:55:22', '2026-01-18 08:55:22'),
('b5adc465-ed24-4525-8e3c-715061ef3238', '1e876554-e267-498d-920e-f09b93d7ab32', 3, '50d3c341-469b-4ea4-bf57-bf7fd5894674', '{\"min_amount\":1000000}', '', '2026-01-18 08:55:34', '2026-01-18 08:55:34');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `actors`
--
ALTER TABLE `actors`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `code` (`code`);

--
-- Indexes for table `approval_history`
--
ALTER TABLE `approval_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_request_id` (`request_id`),
  ADD KEY `idx_actor_id` (`actor_id`),
  ADD KEY `idx_user_id` (`user_id`);

--
-- Indexes for table `requests`
--
ALTER TABLE `requests`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_workflow_id` (`workflow_id`),
  ADD KEY `idx_requester_id` (`requester_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `idx_actor_id` (`actor_id`);

--
-- Indexes for table `workflows`
--
ALTER TABLE `workflows`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `workflow_steps`
--
ALTER TABLE `workflow_steps`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `idx_workflow_level` (`workflow_id`,`level`),
  ADD KEY `idx_workflow_id` (`workflow_id`),
  ADD KEY `idx_actor_id` (`actor_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
