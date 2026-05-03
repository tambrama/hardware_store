from app.models import Address


def address_dao_to_dto(dao: Address) -> dict:
    return {
        "id": str(dao.id),
        "country": dao.country,
        "city": dao.city,
        "street": dao.street,
    }


def address_dto_to_dao(data: dict, dao: Address = None) -> Address:
    if dao is None:
        dao = Address()
    dao.country = data["country"]
    dao.city = data["city"]
    dao.street = data["street"]
    return dao
