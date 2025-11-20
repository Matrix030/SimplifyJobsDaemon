def create_project_selection_prompt(job_description: str, projects: list[dict]) -> str:
    
    projects_text = "\n".join([
        f"{i + 1}. {p['name']}: {p['description']}"
        for i, p in enumerate(projects)
    ])
    

    prompt = f"""You are an expert technical recruiter analyzing job requirements.
    JOB DESCRIPTION:
        {job_description}

    AVAILABLE PROJECTS:
    {projects_text}

    TASK: Select the 2 most relevant projects that best demonstrate skills matching this job description.


    Respond in this exact JSON format:
            {{
            "selected_projects": ["project_name_1", "project_name_2"],
            "reasoning":  "Brief explaination of why these projects match the job requirements"
            }}

    Only return valid JSON, nothing else.
    """

    return prompt


   
