drop table if exists content;
delete type content_type;

drop table if exists blog.entries;
drop table if exists blog.tags;
delete schema blog;

drop table if exists songbook.songs;
delete schema songs;

drop table if exists cv.portfolio;
drop table if exists cv.work_experience;
drop table if exists cv.education;
drop table if exists cv.hobbies;
delete schema cv;

drop table if exists auth.users;
delete schema auth;

delete database kuda_ai;
