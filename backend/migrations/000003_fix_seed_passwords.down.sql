-- Rollback: restore original (broken) hashes
UPDATE users SET password_hash = '$2a$12$9Imjwp7N3Pg9JYrSFXGhBeS/CigEfaTPDumYSInWxuvfAshWNBcMq'
WHERE email = 'teacher@test.local';

UPDATE users SET password_hash = '$2a$12$NpNB3XRy4Oe/6JtkI8oIZecJt8dXe1uxrKsEjiNuO.TUxvTPIBO0O'
WHERE email = 'student@test.local';

UPDATE users SET password_hash = '$2a$12$/87O24Hb8Cysmi69KIfcEebzZHK75PeElK/cDUwjqZiifRPC1Oc0.'
WHERE email = 'qaadmin@test.local';

UPDATE users SET password_hash = '$2a$12$k7.V0YormexAXHDl4WFg/.KLHfgyiOOaj5piW7CTWMGkgUfRpcGga'
WHERE email IN ('qateacher@test.local', 'qastudent@test.local');
