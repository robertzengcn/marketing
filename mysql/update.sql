-- branch scrapyinfo-api
ALTER TABLE mk_social_account DROP COLUMN campaign_id;
ALTER TABLE mk_social_account DROP COLUMN socialplatform;

UPDATE TABLE mk_proxy set account_id=1 where 1=1;

INSERT INTO mk_social_platform (name,url)
	VALUES ('douyin','https://www.douyin.com/');
INSERT INTO amigamarketing.mk_social_platform (name,url)
	VALUES ('tiktok','https://www.tiktok.com/');

ALTER TABLE mk_social_account DROP COLUMN proxy_id;


