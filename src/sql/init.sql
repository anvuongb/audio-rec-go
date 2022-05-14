CREATE TABLE IF NOT EXISTS voice_metadata(
    id integer primary key autoincrement,
    request_id varchar(255),
    file_id varchar(255),
    filepath_mask varchar(255),
    filepath_no_mask varchar(255),
    generated_text varchar(255),
	created_at datetime default CURRENT_TIMESTAMP,
    updated_at datetime default CURRENT_TIMESTAMP,
    unique(request_id),
    unique(file_id)
)