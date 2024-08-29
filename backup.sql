-- Table: public.auth

CREATE TABLE IF NOT EXISTS public.auth
(
    uuid character varying(36) COLLATE pg_catalog."default" NOT NULL,
    username text COLLATE pg_catalog."default" NOT NULL,
    email text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT auth_pkey PRIMARY KEY (uuid),
    CONSTRAINT "uniqEmail" UNIQUE (email),
    CONSTRAINT "uniqUsername" UNIQUE (username)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.auth
    OWNER to postgres;


-- Table: public.refreshtoken

CREATE TABLE IF NOT EXISTS public.refreshtoken
(
    token text COLLATE pg_catalog."default" NOT NULL,
    uuid character varying(36) COLLATE pg_catalog."default" NOT NULL,
    ip text COLLATE pg_catalog."default",
    sessionid integer,
    CONSTRAINT pk PRIMARY KEY (uuid),
    CONSTRAINT "uniqSession" UNIQUE (sessionid),
    CONSTRAINT "uniqUUID" UNIQUE (uuid),
    CONSTRAINT "fkUUID" FOREIGN KEY (uuid)
        REFERENCES public.auth (uuid) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.refreshtoken
    OWNER to postgres;


-- FUNCTION: public.getrefresh(character varying)

CREATE OR REPLACE FUNCTION public.getrefresh(
	uuid character varying)
    RETURNS void
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
      BEGIN
       SELECT token from refreshToken where uuidForeign = uuid;
      END;
  
$BODY$;

ALTER FUNCTION public.getrefresh(character varying)
    OWNER TO postgres;


-- FUNCTION: public.register(text, text, text)

CREATE OR REPLACE FUNCTION public.register(
	username text,
	email text,
	uuid text)
    RETURNS void
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
      BEGIN
        INSERT INTO auth(username, email, uuid)
        VALUES(username, email, uuid);
      END;
  
$BODY$;

ALTER FUNCTION public.register(text, text, text)
    OWNER TO postgres;


-- FUNCTION: public.searchemail(character varying)

CREATE OR REPLACE FUNCTION public.searchemail(
	uuidvalue character varying)
    RETURNS SETOF text 
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
      BEGIN
      return query  SELECT email from auth where uuid = uuidValue limit 1;
      END;
  
$BODY$;

ALTER FUNCTION public.searchemail(character varying)
    OWNER TO postgres;


-- FUNCTION: public.searchrefresh(text, integer)

CREATE OR REPLACE FUNCTION public.searchrefresh(
	ipaddress text,
	session integer)
    RETURNS SETOF text 
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
      BEGIN
      return query SELECT token from refreshToken where ip = ipaddress and sessionid = session limit 1;
      END;
  
$BODY$;

ALTER FUNCTION public.searchrefresh(text, integer)
    OWNER TO postgres;


-- FUNCTION: public.setrefreshtoken(character varying, text, text, integer)

CREATE OR REPLACE FUNCTION public.setrefreshtoken(
	uuidnew character varying,
	refreshtokenvalue text,
	ipaddress text,
	sessionid integer)
    RETURNS void
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
      BEGIN
        INSERT INTO refreshToken(token, uuid,ip,sessionid)
        VALUES(refreshTokenValue, uuidNew,ipaddress,sessionid);
      END;
  
$BODY$;

ALTER FUNCTION public.setrefreshtoken(character varying, text, text, integer)
    OWNER TO postgres;


-- FUNCTION: public.updaterefreshtoken(character varying, text, text, integer)

CREATE OR REPLACE FUNCTION public.updaterefreshtoken(
	uuidnew character varying,
	refreshtokenvalue text,
	ipaddress text,
	sessionidvalue integer)
    RETURNS void
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
      BEGIN
	  UPDATE refreshtoken
	SET token=refreshTokenValue, uuid=uuidNew, ip=ipaddress, sessionid=sessionidvalue
	WHERE uuid=uuidNew;
      END;
  
$BODY$;

ALTER FUNCTION public.updaterefreshtoken(character varying, text, text, integer)
    OWNER TO postgres;
