
DROP TABLE IF EXISTS `person`;
CREATE TABLE `person` (
  `id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(500) DEFAULT NULL,
  `full_name` varchar(500) DEFAULT NULL,
  `email` varchar(90) DEFAULT NULL,
  `phone_number` varchar(20) DEFAULT NULL,
  `address_city` varchar(500) DEFAULT NULL,
  `address_street` varchar(500) DEFAULT NULL,
  `address_house` varchar(500) DEFAULT NULL,
  `address_entrance` varchar(500) DEFAULT NULL,
  `address_floor` varchar(500) DEFAULT NULL,
  `address_office` varchar(500) DEFAULT NULL,
  `address_comment` varchar(500) DEFAULT NULL,
  `location_latitude` varchar(20) DEFAULT NULL,
  `location_longitude` varchar(20) DEFAULT NULL,
  `amount_charged` varchar(20) DEFAULT NULL,
  `user_id` varchar(100) DEFAULT NULL,
  `user_agent` varchar(1000) DEFAULT NULL,
  `created_at` varchar(50) DEFAULT NULL,
  `address_doorcode` varchar(300) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `age` (`phone_number`),
  KEY `num` (`phone_number`)
) ENGINE=InnoDB AUTO_INCREMENT=30866986 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
LOCK TABLES `person` WRITE;
INSERT INTO `person` VALUES (1,'Арина','','','+75555555555','Екатеринбург','Московская улица','75','','17','247','','56.827372','60.586806','608','16765503','Mozilla/5.0 (iPhone; CPU iPhone OS 12_1_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) yandex-taxi/650.12.0.165399 YandexEatsKit/1.26.0 Superapp/Eats','2021-10-21T11:15:38.000Z','+');
UNLOCK TABLES;
