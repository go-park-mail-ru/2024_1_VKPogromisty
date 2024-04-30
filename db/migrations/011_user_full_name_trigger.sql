-- Write your migrate up statements here
-- Create the trigger function
CREATE OR REPLACE FUNCTION update_full_name()
RETURNS TRIGGER AS $$
BEGIN
    NEW.full_name := NEW.first_name || ' ' || NEW.last_name;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER update_full_name_trigger
BEFORE INSERT OR UPDATE ON public.user
FOR EACH ROW EXECUTE PROCEDURE update_full_name();
---- create above / drop below ----
DROP TRIGGER IF EXISTS update_full_name_trigger ON public.user;
DROP FUNCTION IF EXISTS update_full_name();
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
