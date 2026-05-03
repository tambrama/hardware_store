from app.models import Supplier
from app.mappers.address_mapper import address_dao_to_dto


def supplier_dao_to_dto(dao: Supplier) -> dict:
    return {
        "id": str(dao.id),
        "name": dao.name,
        "address_id": str(dao.address_id),
        "phone_number": dao.phone_number,
        "address": address_dao_to_dto(dao.address) if dao.address else None,
    }


def supplier_dto_to_dao(data: dict, dao: Supplier = None) -> Supplier:
    if dao is None:
        dao = Supplier()
    dao.name = data["name"]
    dao.address_id = data["address_id"]
    dao.phone_number = data["phone_number"]
    return dao
