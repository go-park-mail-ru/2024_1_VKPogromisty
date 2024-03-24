-- Write your migrate up statements here
CREATE OR REPLACE FUNCTION trigger_set_timestamp() 
RETURNS TRIGGER AS $$ 
BEGIN 
  NEW.updated_at = NOW(); 
  RETURN NEW; 
  END; 
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.user 
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.post
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.personal_message 
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.subscription
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.post_attachment
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();
---- create above / drop below ----
DROP TRIGGER IF EXISTS set_timestamp ON public.user;
DROP TRIGGER IF EXISTS set_timestamp ON public.post;
DROP TRIGGER IF EXISTS set_timestamp ON public.personal_message;
DROP TRIGGER IF EXISTS set_timestamp ON public.subscription;
DROP TRIGGER IF EXISTS set_timestamp ON public.post_attachment;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
