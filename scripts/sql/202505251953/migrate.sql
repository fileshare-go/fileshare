CREATE TABLE `file_infos` (
  `filename` text NOT NULL,
  `sha256` text NOT NULL,
  `chunk_size` bigint,
  `chunk_number` integer,
  `file_size` bigint,
  `uploaded_chunks` text,
  PRIMARY KEY (`sha256`)
);

CREATE TABLE `share_links` (
  `link_code` text,
  `sha256` text,
  PRIMARY KEY (`sha256`),
  CONSTRAINT `fk_file_infos_link`
  FOREIGN KEY (`sha256`)
  REFERENCES `file_infos`(`sha256`)
);