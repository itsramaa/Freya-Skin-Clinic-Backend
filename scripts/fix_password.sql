-- Quick fix: Reset password ke hash yang benar untuk "admin"
UPDATE users SET 
    password_hash = '$2a$10$deA8Pp1xEpWu1iTjTy6O0.HI/cmdnsvE.2VYEBXSUBSOTYWOkCBuC',
    is_default_password = true,
    session_id = NULL,
    updated_at = NOW();

-- Verifikasi
SELECT username, is_default_password FROM users;
