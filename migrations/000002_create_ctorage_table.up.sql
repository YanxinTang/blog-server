CREATE TABLE public.storage (
    id          bigserial NOT NULL,
    name        text,
    secret_id   text,
    secret_key  text,
    token       text,
    region      text,
    endpoint    text,
    bucket      text,
    usage       bigint DEFAULT 0,
    capacity    bigint,
    updated_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT storage_pkey PRIMARY KEY (id)
);

CREATE TRIGGER update_storage BEFORE UPDATE ON public.storage FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
