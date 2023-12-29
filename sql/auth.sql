CREATE OR REPLACE FUNCTION register(
    fullname TEXT,
    p_is_admin BOOLEAN,
    p_password TEXT,
    p_email TEXT,
    p_phone TEXT,
    p_street TEXT,
    p_house TEXT
) RETURNS TEXT
    LANGUAGE plpgsql AS
$$

DECLARE
    new_owner_id INTEGER := 0;
    owner_house_id INTEGER;
    password_hash TEXT;
BEGIN
    SELECT House_ID
    FROM House WHERE House.Street = p_street AND House.Number_of_House = p_house INTO owner_house_id;
    IF owner_house_id IS NULL THEN
        RETURN 'House not found';
    END IF ;

    SELECT User_ID
    FROM "Users"
    WHERE Full_Name = fullname
    INTO new_owner_id;

    IF new_owner_id IS NOT NULL THEN
        RETURN 'User already exists';
    END IF;

    SELECT md5(p_password)
    INTO password_hash;

    INSERT INTO "Users" (Full_Name, Is_Admin, "Password", House_ID, Email, Phone_Number)
    VALUES (fullname, p_is_admin, password_hash, owner_house_id, p_email, p_phone)
    RETURNING "Users".User_ID INTO new_owner_id;
    RETURN CONCAT(CAST(new_owner_id AS TEXT), ' ', CAST(p_is_admin AS TEXT));
END;
$$;

CREATE OR REPLACE FUNCTION LOGIN(fullname text, password text)
    RETURNS TEXT
    LANGUAGE plpgsql AS
$$
DECLARE
    expected_password TEXT;
    received_password TEXT;
    current_user_id INTEGER;
    current_admin BOOLEAN;
BEGIN
    SELECT "Users"."Password", "Users".User_ID, "Users".Is_Admin
    FROM "Users"
    WHERE "Users".Full_Name = login.fullname
    INTO expected_password, current_user_id, current_admin;

    IF NOT FOUND THEN
        RETURN 'Invalid fullname';
    END IF;

    SELECT md5(password)
    INTO received_password;

    IF expected_password <> received_password THEN
        RETURN 'Invalid password';
    END IF;

    RETURN CONCAT(CAST(current_user_id AS TEXT), ' ', CAST(current_admin AS TEXT));

END;
$$;