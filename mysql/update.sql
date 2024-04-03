-- branch scrapyinfo-api
ALTER TABLE mk_social_account DROP COLUMN campaign_id;
ALTER TABLE mk_social_account DROP COLUMN socialplatform;

UPDATE TABLE mk_proxy set account_id=1 where 1=1;

