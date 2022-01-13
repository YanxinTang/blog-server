CREATE TABLE public."user" (
    id bigserial NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    salt text NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT user_unique_username UNIQUE (username),
    CONSTRAINT user_unique_email UNIQUE (email)
);


CREATE TABLE public.category (
    id bigserial NOT NULL,
    name text,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT category_pkey PRIMARY KEY (id),
    CONSTRAINT category_unique_name UNIQUE (name)
);


CREATE TABLE public.article (
    id bigserial NOT NULL,
    category_id bigint,
    title text,
    content text,
    status boolean,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT article_pkey PRIMARY KEY (id),
    CONSTRAINT article_fk_category_id FOREIGN KEY (category_id) REFERENCES public.category(id) ON UPDATE RESTRICT ON DELETE CASCADE
);


CREATE TABLE public.comment (
    id bigserial NOT NULL,
    article_id bigint NOT NULL,
    parent_id bigint DEFAULT 0 NOT NULL,
    username text NOT NULL,
    content text NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT comment_pkey PRIMARY KEY (id),
    CONSTRAINT comment_fk_article_id FOREIGN KEY (article_id) REFERENCES public.article(id) ON UPDATE RESTRICT ON DELETE CASCADE
);


CREATE TABLE public.setting (
    id bigserial NOT NULL,
    key text NOT NULL,
    value text NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT setting_pkey PRIMARY KEY (id),
    CONSTRAINT setting_unique_key UNIQUE (key)
);


CREATE FUNCTION public.set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
   IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.updated_at = now(); 
      RETURN NEW;
   ELSE
      RETURN OLD;
   END IF;
END;
$$;


CREATE TRIGGER update_article BEFORE UPDATE ON public.article FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER update_category BEFORE UPDATE ON public.category FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER update_comment BEFORE UPDATE ON public.comment FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER update_setting BEFORE UPDATE ON public.setting FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
CREATE TRIGGER update_user BEFORE UPDATE ON public."user" FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


INSERT INTO setting (key, value) VALUES ('signupEnable', 1);