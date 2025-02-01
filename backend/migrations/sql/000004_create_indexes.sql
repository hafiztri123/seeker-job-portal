CREATE INDEX idx_users_email ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX idx_companies_email ON companies (email) WHERE deleted_at IS NULL;
CREATE INDEX idx_jobs_company ON jobs (company_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_jobs_status ON jobs (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_jobs_location_type ON jobs (location_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_jobs_deadline ON jobs (deadline) WHERE deleted_at IS NULL AND status = 'open';

CREATE INDEX idx_jobs_search ON jobs USING gin(
    to_tsvector('english',
        COALESCE(title,'') || ' ' ||
        COALESCE(description,'') || ' '
    )
) WHERE deleted_at IS NULL;

CREATE INDEX idx_jobs_filters ON jobs (
    location_type,
    employment_type,
    status
) WHERE deleted_at IS NULL;

CREATE INDEX idx_job_applications_status ON job_applications (user_id, status);
CREATE INDEX idx_user_skills_level ON user_skills (user_id, proficiency_level);
CREATE INDEX idx_job_skills_required ON job_skills (job_id, required_level);