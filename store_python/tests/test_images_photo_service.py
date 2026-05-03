import uuid
from decimal import Decimal
from unittest.mock import patch

from app.extensions import db
from app.models import Address, Supplier, Product


class _FakeRedis:
    def __init__(self):
        self.kv = {}

    def get(self, k):
        return self.kv.get(k)

    def setex(self, k, ttl, v):
        self.kv[k] = v

    def delete(self, k):
        self.kv.pop(k, None)


def _mk_product():
    addr = Address(country="RU", city="Moscow", street="Tverskaya")
    db.session.add(addr)
    db.session.commit()

    supp = Supplier(name="S1", address_id=addr.id, phone_number="+70000000000")
    db.session.add(supp)
    db.session.commit()

    p = Product(
        name="P1",
        category="C1",
        price=Decimal("10.00"),
        available_stock=1,
        supplier_id=supp.id,
    )
    db.session.add(p)
    db.session.commit()
    return p


def test_add_image_calls_photo_service_and_sets_product_image_id(client, auth_headers):
    product = _mk_product()
    fake_id = str(uuid.uuid4())

    class _Resp:
        status_code = 200

        def json(self):
            return {"id": fake_id}

    with patch("app.api.images.requests.post", return_value=_Resp()):
        resp = client.post(
            f"/api/v1/images?product_id={product.id}",
            data=b"\x89PNG\r\n\x1a\nxxxx",
            headers={**auth_headers, "Content-Type": "application/octet-stream"},
        )
    assert resp.status_code == 201

    updated = db.session.get(Product, product.id)
    assert str(updated.image_id) == fake_id


def test_get_image_by_product_id_caches_response(client, auth_headers):
    product = _mk_product()
    product.image_id = uuid.uuid4()
    db.session.commit()

    img_bytes = b"binary-image"

    class _Resp:
        status_code = 200
        content = img_bytes

    fake_redis = _FakeRedis()

    with patch("app.api.images.get_redis_client", return_value=fake_redis):
        with patch("app.api.images.requests.get", return_value=_Resp()):
            resp1 = client.get(f"/api/v1/images/product/{product.id}", headers=auth_headers)
            assert resp1.status_code == 200
            assert resp1.data == img_bytes

            # second call should HIT cache and not call requests.get again
            with patch("app.api.images.requests.get") as mocked_get:
                resp2 = client.get(f"/api/v1/images/product/{product.id}", headers=auth_headers)
                assert resp2.status_code == 200
                assert resp2.data == img_bytes
                mocked_get.assert_not_called()

