package queryfileupload

const (
	QuerySave = `
	INSERT INTO public.sa_file_upload (file_name, file_path, file_type, user_input, user_edit) 
	VALUES(:file_name, :file_path, :file_type, :user_input, :user_edit);

	`
)
