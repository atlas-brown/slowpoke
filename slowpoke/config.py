def get_request_ratio(benchmark, request="mix"):
    if benchmark == "boutique":
        if request == "mix":
            return {
                'cart': 0.45175405908969923, 
                'checkout': 0.05106201756720788, 
                'currency': 0.2597231833910035, 
                'email': 0.0, 
                'frontend': 1.0, 
                'payment': 0.05106201756720788, 
                'product_catalog': 0.6091136545115784, 
                'recommendations': 0.0, 
                'shipping': 0.10212403513441576
            }
        elif request == "home":
            return {
                'cart': 1.0, 
                'checkout': 1.0, 
                'currency': 1.0, 
                'email': 1.0, 
                'frontend': 1.0, 
                'payment': 1.0, 
                'product_catalog': 1.0, 
                'recommendations': 1.0, 
                'shipping': 1.0
            }
        else:
            raise ValueError(f"[config.py] Unknown request type: {request} for benchmark: {benchmark}")
    elif benchmark == "social":
        return {
            "composepost": 0.09979209979,
            "hometimeline": 0.7005428505,
            "poststorage": 0.9497747748,
            "socialgraph": 0.099792099790,
            "usertimeline": 0.4022002772
        }
    elif benchmark == "movie":
        return {
            "castinfo": 0.9050900514,
            "composereview": 0.1003336848,
            "frontend": 0.1003336848,
            "movieid": 0.1003336848,
            "movieinfo": 0.9050900514,
            "moviereviews": 0.9050900514,
            "page": 0.9050900514,
            "plot": 0.9050900514,
            "reviewstorage": 1.005423736183,
            "uniqueid": 0.1003336848,
            "user": 0.2006673695,
            "userreviews": 0.1003336848
        }
    elif benchmark == "hotel":
        return {
            "frontend": 1,
            "profile": 0.8,
            "rate": 0.8,
            "reservation": 1,
            "search": 0.8,
            "user": 0.198478617
        }
    else:
        raise ValueError(f"[config.py] Unknown benchmark: {benchmark}")

def get_baseline_service_processing_time(benchmark, request="mix"):
    if benchmark == "boutique":
        return {
            "cart":0,
            "checkout":0,
            "currency":0,
            "email":0,
            "frontend":0,
            "payment":0,
            "product_catalog":0,
            "recommendations":0,
            "shipping":0
        }
    elif benchmark == "social":
        return {
            "composepost": 0,
            "hometimeline": 0,
            "poststorage": 0,
            "socialgraph": 0,
            "usertimeline": 0
        }
    elif benchmark == "movie":
        return {
            "castinfo": 0,
            "composereview": 0,
            "frontend": 0,
            "movieid": 0,
            "movieinfo": 0,
            "moviereviews": 0,
            "page": 0,
            "plot": 0,
            "reviewstorage": 0,
            "uniqueid": 0,
            "user": 0,
            "userreviews": 0
        }
    elif benchmark == "hotel":
        return {
            "frontend": 0,
            "profile": 0,
            "rate": 0,
            "reservation": 0,
            "search": 0,
            "user": 0
        }
    else :
        raise ValueError(f"[config.py] Unknown benchmark: {benchmark}")