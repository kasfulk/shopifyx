/*
 Navicat Premium Data Transfer

 Source Server         : PC
 Source Server Type    : PostgreSQL
 Source Server Version : 150005 (150005)
 Source Host           : 100.66.78.33:5432
 Source Catalog        : shopifyx
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 150005 (150005)
 File Encoding         : 65001

 Date: 11/03/2024 10:53:27
*/


-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" int8 NOT NULL GENERATED ALWAYS AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1
),
  "username" varchar(255) COLLATE "pg_catalog"."default",
  "password" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" date DEFAULT now(),
  "updated_at" date,
  "deleted_at" date
)
;
ALTER TABLE "public"."users" OWNER TO "postgres";

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");
