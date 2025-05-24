CREATE TABLE `file_infos` (
  `filename` text NOT NULL,
  `sha256` text NOT NULL,
  `chunk_size` bigint,
  `chunk_number` integer,
  `file_size` bigint,
  `uploaded_chunks` text,
  PRIMARY KEY (`sha256`)
);
