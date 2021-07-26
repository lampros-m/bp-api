
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `bestprice` /*!40100 DEFAULT CHARACTER SET latin1 */;

USE `bestprice`;
DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `categories` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `place` int(10) unsigned NOT NULL,
  `url_image` varchar(2083) COLLATE utf8_unicode_ci NOT NULL,
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_on` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `place` (`place`),
  KEY `categories_place_IDX` (`place`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='BestPrice categories';
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'house',1,'http://www.archive.bp.gr/images/categories/house.webp','2020-11-01 12:53:43',NULL),(2,'garden',2,'http://www.archive.bp.gr/images/categories/garden.webp','2020-11-01 12:53:43',NULL),(3,'technology',3,'http://www.archive.bp.gr/images/categories/technology.webp','2020-11-01 12:53:43',NULL),(4,'kids',4,'http://www.archive.bp.gr/images/categories/kids.webp','2020-11-01 12:53:43',NULL),(5,'sports',5,'http://www.archive.bp.gr/images/categories/sports.webp','2020-11-01 12:53:43',NULL);
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `products` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `category_id` int(10) unsigned DEFAULT NULL,
  `title` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `url_image` varchar(2083) COLLATE utf8_unicode_ci NOT NULL,
  `price` decimal(6,2) NOT NULL,
  `description` text COLLATE utf8_unicode_ci,
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_on` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_products_categories` (`category_id`),
  CONSTRAINT `FK_products_categories` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='BestPrice products';
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES (1,1,'fridge','http://www.archive.bp.gr/images/products/fridge.webp',80.40,'this is a fridge','2020-11-01 12:53:43',NULL),(2,1,'bed','http://www.archive.bp.gr/images/products/bed.webp',61.20,'this is a bed','2020-11-01 12:53:43',NULL),(3,1,'chair','http://www.archive.bp.gr/images/products/chair.webp',15.35,'this is a chair','2020-11-01 12:53:43',NULL),(4,2,'grill','http://www.archive.bp.gr/images/products/grill.webp',37.70,'this is a grill','2020-11-01 12:53:43',NULL),(5,2,'fence','http://www.archive.bp.gr/images/products/fence.webp',22.40,'this is a fence','2020-11-01 12:53:43',NULL),(6,2,'flowerpot','http://www.archive.bp.gr/images/products/flowerpot.webp',4.80,'this is a flowerpot','2020-11-01 12:53:43',NULL),(7,3,'cellphone','http://www.archive.bp.gr/images/products/cellphone.webp',252.40,'this is a cellphone','2020-11-01 12:53:43',NULL),(8,3,'laptop','http://www.archive.bp.gr/images/products/laptop.webp',751.20,'this is a laptop','2020-11-01 12:53:43',NULL),(9,3,'rasberry-pi','http://www.archive.bp.gr/images/products/rasberry-pi.webp',65.40,'this is a rasberry pi','2020-11-01 12:53:43',NULL),(10,4,'doll','http://www.archive.bp.gr/images/products/doll.webp',12.37,'this is a doll','2020-11-01 12:53:43',NULL),(11,4,'miniature','http://www.archive.bp.gr/images/products/miniature.webp',9.99,'this is a miniature','2020-11-01 12:53:43',NULL),(12,4,'puzzle','http://www.archive.bp.gr/images/products/puzzle.webp',12.37,'this is a puzzle','2020-11-01 12:53:43',NULL),(13,5,'ball','http://www.archive.bp.gr/images/products/ball.webp',6.27,'this is a ball','2020-11-01 12:53:43',NULL),(14,5,'dumbbell','http://www.archive.bp.gr/images/products/dumbbell.webp',8.88,'this is a dumbbell','2020-11-01 12:53:43',NULL),(15,5,'sneaker','http://www.archive.bp.gr/images/products/sneaker.webp',34.17,'this is a sneaker','2020-11-01 12:53:43',NULL);
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `username` varchar(254) NOT NULL,
  `pass` text NOT NULL,
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_username_key` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('e9e312a0-1c5e-11eb-999f-0242ac160002','johndoe','$2a$08$D67CAZ2OkrClttchwarNtucWTOsHDF3r2bnV/s3Q2oTwm8oTqb/ni','2020-11-01 16:25:54');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

