package queryauth

const (
	QueryAuthLogin = `
	select a.user_id ,
		a.name ,
		a.email ,
		a.telp ,
		a.join_date ,
		a.user_type ,
		a.file_id ,
		b.file_name ,
		b.file_path ,
		a.pwd
	from public.ss_user a
	left join sa_file_upload b
		on a.file_id = b.file_id 
	where (a.email ilike $1 or a.telp = $2);    `

	QueryUpdatePassword = `
	UPDATE public.ss_user
		SET pwd = :pwd
	WHERE user_id = :user_id;
	`
)
