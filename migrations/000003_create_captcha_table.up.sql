CREATE TABLE public.captcha (
    id          bigserial NOT NULL,
    key         text UNIQUE,
    text        text,
    updated_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT captcha_pkey PRIMARY KEY (id)
);

CREATE TRIGGER update_captcha BEFORE UPDATE ON public.storage FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
