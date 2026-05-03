from datetime import datetime
from app.models import Product


def _parse_date(v):
    if v is None:
        return None
    if hasattr(v, "isoformat"):
        return v
    if isinstance(v, str):
        return datetime.fromisoformat(v.replace("Z", "+00:00")).date()
    return v


def product_dao_to_dto(dao: Product) -> dict:
    return {
        "id": str(dao.id),
        "name": dao.name,
        "category": dao.category,
        "price": str(dao.price) if dao.price is not None else None,
        "available_stock": dao.available_stock,
        "last_update_date": dao.last_update_date.isoformat() if dao.last_update_date else None,
        "supplier_id": str(dao.supplier_id),
        "image_id": str(dao.image_id) if dao.image_id else None,
    }


def product_dto_to_dao(data: dict, dao: Product = None) -> Product:
    if dao is None:
        dao = Product()
    dao.name = data["name"]
    dao.category = data["category"]
    dao.price = data["price"]
    dao.available_stock = data.get("available_stock", 0)
    dao.last_update_date = _parse_date(data.get("last_update_date"))
    dao.supplier_id = data["supplier_id"]
    dao.image_id = data.get("image_id")
    return dao
