package queryoption

const (
	QueryGetListOption = `SELECT 	option_id, option_url, method_api, 
									sp, line_no, table_name, 
									user_input, user_edit, time_input, 
									time_edit
						  FROM public.ss_option_db
						  		WHERE option_url iLIKE $1
						  ORDER BY line_no , method_api;					  
	`
	QueryGetListParamFunction = ` 
			SELECT routines.routine_name,
				parameters.parameter_name,
				parameters.data_type,
				parameters.ordinal_position
			FROM information_schema.routines
			LEFT JOIN information_schema.parameters ON routines.specific_name = parameters.specific_name
			WHERE 	routines.specific_catalog = current_database() AND
					routines.specific_schema = 'public' AND
					routines.routine_name iLIKE $1 AND
					parameters.parameter_mode = 'IN'
			ORDER BY routines.specific_name, routines.routine_name, parameters.ordinal_position;

	
	`

	// SELECT routine_name, parameter_name, data_type, oridinal_position
	// 							 FROM public.get_param_function($1) order by oridinal_position asc;
	QueryGetListFieldType = `
				select 	ordinal_position,
						column_name as parameter_name,
						data_type,
						case when character_maximum_length is not null
						then character_maximum_length
						else numeric_precision end as max_length,
						is_nullable,
						numeric_precision as precision, numeric_scale as scale,
						column_default as default_value
				from information_schema.columns
				where table_schema not in ('information_schema', 'pg_catalog')
				AND table_name iLIKE $1
				ORDER BY ordinal_position;
	`

	QueryResultFunctionType = `
			SELECT routines.routine_name,
				parameters.parameter_name,
				parameters.data_type,
				parameters.ordinal_position
			FROM information_schema.routines
			LEFT JOIN information_schema.parameters ON routines.specific_name = parameters.specific_name
			WHERE 	routines.specific_catalog = current_database() AND
					routines.specific_schema = 'public' AND
					routines.routine_name iLIKE $1 AND
					parameters.parameter_mode <> 'IN'
			ORDER BY routines.specific_name, routines.routine_name, parameters.ordinal_position;
	`

	QueryExecCUD = `
		SELECT {FunctionName} as row_id FROM public.{FunctionName} ({ParameterFunction});
	`
	QueryGetByID = `
		SELECT * FROM public.{FunctionName} ({ParameterFunction});
	`
	QueryList = `
		WITH result_set AS
		(
			SELECT
				row_number() OVER ({sSortFiled}) as no,
				{sField}
			FROM {sTable}
			{sWhere}
		)
		SELECT * FROM result_set
		LIMIT $1
		OFFSET $2 ;
	`

	QueryDefineColumn = `
			SELECT column_field
			FROM ss_define_column
			WHERE option_url iLIKE $1 AND line_no = $2;
	`
)
