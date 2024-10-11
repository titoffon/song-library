INSERT INTO songs (music_group, song, releasedate, text, link)
VALUES
    ('Muse', 'Supermassive Black Hole', '2006-06-19', E'Oh baby don''t \n\n you know I suffer?', 'https://example.com/supermassive-black-hole'),
    ('Radiohead', 'Creep', '1993-09-21', E'I''m a creep,\n\n I''m a weirdo', 'https://example.com/creep'),
    ('Coldplay', 'Fix You', '2005-09-05', E'Lights will\n\n guide you home', 'https://example.com/fix-you'),
    ('Muse', 'Uprising', '2009-09-07', E'They will not force us,\n\n they will stop degrading us', 'https://example.com/uprising'),
    ('The Beatles', 'Hey Jude', '1968-08-26', E'Hey Jude, don''t\n\n make it bad', 'https://example.com/hey-jude'),
    ('Pink Floyd', 'Comfortably Numb', '1979-11-30', E'Hello, is there\n\n anybody in there?', 'https://example.com/comfortably-numb'),
    ('Radiohead', 'Karma Police', '1997-08-25', E'This is what you''ll\n\n get, when you mess with us', 'https://example.com/karma-police'),
    ('Coldplay', 'Viva La Vida', '2008-05-25', E'I used to rule the world,\n\n seas would rise when I gave the word', 'https://example.com/viva-la-vida'),
    ('The Beatles', 'Let It Be', '1970-03-06', E'When I find myself in times of trouble,\n\n Mother Mary comes to me', 'https://example.com/let-it-be'),
    ('Nirvana', 'Smells Like Teen Spirit', '1991-09-10', E'With the lights out,\n\n it''s less dangerous', 'https://example.com/smells-like-teen-spirit'),
    ('Queen', 'Bohemian Rhapsody', '1975-10-31', E'Is this the real life?\n\n Is this just fantasy?', 'https://example.com/bohemian-rhapsody'),
    ('Muse', 'Madness', '2012-08-20', E'I have finally\n\n seen the light', 'https://example.com/madness'),
    ('Pink Floyd', 'Wish You Were Here', '1975-09-12', E'How I wish,\n\n how I wish you were here', 'https://example.com/wish-you-were-here'),
    ('Radiohead', 'No Surprises', '1998-01-12', E'A handshake of carbon\n\n monoxide', 'https://example.com/no-surprises'),
    ('The Beatles', 'Come Together', '1969-10-06', E'Here come old flat\n\n top, he come groovin'' up slowly', 'https://example.com/come-together'),
    ('Coldplay', 'Paradise', '2011-09-12', E'When she was just a girl, she \n\n expected the world', 'https://example.com/paradise'),
    ('Queen', 'Another One Bites the Dust', '1980-08-22', E'Are you ready?\n\n Hey, are you ready for this?', 'https://example.com/another-one-bites-the-dust'),
    ('Nirvana', 'Lithium', '1991-09-13', E'I''m so happy, cause today\n\n I found my friends', 'https://example.com/lithium'),
    ('The Beatles', 'Yellow Submarine', '1966-08-05', E'We all live in\n\n a yellow submarine', 'https://example.com/yellow-submarine'),
    ('Pink Floyd', 'Money', '1973-03-01', E'Money,\n\n get away', 'https://example.com/money')
ON CONFLICT (music_group, song) DO NOTHING;
