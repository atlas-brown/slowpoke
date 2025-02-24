def get_request_ratio(benchmark):
    if benchmark == "boutique":
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
    else:
        raise ValueError(f"[config.py] Unknown benchmark: {benchmark}")

def get_baseline_service_processing_time(benchmark):
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
    else :
        raise ValueError(f"[config.py] Unknown benchmark: {benchmark}")