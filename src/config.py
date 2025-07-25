#!/usr/bin/env python3

import random

service_reuse = {
    "dag-relay": [1, 1, 1, 1, 3, 3, 3],
    "dag-cross": [1, 1, 2, 1, 4],
    "dynamic-once": [1, 1, 0.8, 0.2, 0.8],
    "dynamic-twice": [1, 0.2, 0.3, 0.5, 0.16, 0.15, 0.35],
    "dynamic-cache": [1, 1, 0.25],
    "dynamic-cycle": [1, 1.34, 0.34, 0.5, 0.5],
    "chain-d2": [1, 1, 1],
    "chain-d4": [1, 1, 1, 1, 1, 1],
    "chain-d8": [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    "fanout-w3": [1, 1, 1, 1],
    "fanout-w5": [1, 1, 1, 1, 1, 1],
    "fanout-w7": [1, 1, 1, 1, 1, 1, 1, 1],
    "fanout-w10": [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    "dag-balanced": [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    "dag-unbalanced": [1, 1, 1, 1, 1, 1, 1, 1]
}

leaf_nodes = {
    "chain-d2": [2],
    "chain-d4": [5],
    "chain-d8": [9],
    "fanout-w3": [1,2,3],
    "fanout-w10": [1,2,3,4,5,6,7,8,9,10],
    "fanout-w7": [1,2,3,4,5,6,7],
    "fanout-w5": [1,2,3,4,5],
    "dag-balanced": [4,5,6,7,8,9,10,11,12], 
    "dag-unbalanced": [5,6],
    "dag-relay": [5,6],
    "dag-cross": [4],
    "dynamic-once": [3,4],
    "dynamic-cache": [2],
    "dynamic-cycle": [3,4],
    "dynamic-twice": [2,5,6]
}

def get_request_ratio(benchmark, request):
    if benchmark == "synthetic":
        return get_request_ratio_synthetic(request)
    elif benchmark == "boutique":
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

def get_request_ratio_synthetic(request):
    topology = "-".join(request.split("-")[:2])
    ratio_arr = service_reuse[topology]
    ratio = {}
    for i in range(len(ratio_arr)):
        ratio[f"service{i}"] = ratio_arr[i]
    return ratio

def get_baseline_service_processing_time(benchmark, request, target, random_seed):
    if benchmark == "synthetic":
        return get_baseline_service_processing_time_synthetic(target, request, random_seed)
    elif benchmark == "boutique":
        p_t = {
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
        p_t[target] = 1000
        return p_t
    elif benchmark == "social":
        p_t = {
            "composepost": 0,
            "hometimeline": 0,
            "poststorage": 0,
            "socialgraph": 0,
            "usertimeline": 0
        }
        p_t[target] = 1000
        return p_t
    elif benchmark == "movie":
        p_t = {
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
        p_t[target] = 1000
        return p_t
    elif benchmark == "hotel":
        p_t = {
            "frontend": 0,
            "profile": 0,
            "rate": 0,
            "reservation": 0,
            "search": 0,
            "user": 0
        }
        p_t[target] = 1000
        return p_t
    else :
        raise ValueError(f"[config.py] Unknown benchmark: {benchmark}")

# def get_baseline_service_processing_time_synthetic(target, request, random_seed):
#     topology = "-".join(request.split("-")[:2])
#     num = len(service_reuse[topology])
#     random.seed(random_seed)
#     random_numbers = [random.gauss(700, 300) for i in range(num)]
#     random_numbers = [abs(r) for r in random_numbers]   # just in case
#     random_numbers.sort()
#     print(f"[config.py] Random numbers for execution time: {random_numbers}")
#     # use the random number to also decide if we want the target service is the bottleneck
#     is_bottleneck = random.choice([True,True,False])
#     print(f"[config.py] Is bottleneck: {is_bottleneck}")
#     picked_service = target
#     if not is_bottleneck:
#         # pick a random service that is not the target service
#         while picked_service == target:
#             picked_service = f"service{random.randint(0, num-1)}"
#     processing_time = {}
#     picked_idx = int(picked_service.replace("service", ""))
#     processing_time[picked_service] = round(2 * random_numbers.pop(-1)/ service_reuse[topology][picked_idx], 2)

#     for i in range(num):
#         service_name = f"service{i}"
#         if service_name == picked_service:
#             continue
#         processing_time[f"service{i}"] = round(random_numbers.pop(-1) / service_reuse[topology][i], 2)
#     return processing_time

# def get_baseline_service_processing_time_synthetic(target, request, random_seed):
#     topology = "-".join(request.split("-")[:2])
#     num = len(service_reuse[topology])
#     random.seed(random_seed)
#     random_numbers = [random.gauss(500, 300) for i in range(num)]
#     random_numbers = [abs(r) + 200 for r in random_numbers]   # just in case
#     random_numbers.sort()
#     print(f"[config.py] Random numbers for execution time: {random_numbers}")

#     processing_time = {}
#     for i in range(num):
#         if i in leaf_nodes[topology] and f"service{i}" == target:
#             processing_time[f"service{i}"] = round(2 * random_numbers.pop(-1) * max(service_reuse[topology]) / service_reuse[topology][i], 2)
#         elif f"service{i}" == target:
#             processing_time[f"service{i}"] = round(random_numbers.pop(-1) * max(service_reuse[topology]) / service_reuse[topology][i], 2)
#         else:
#             processing_time[f"service{i}"] = round(random_numbers.pop(-1) * max(service_reuse[topology]) / service_reuse[topology][i], 2)
#     return processing_time

def get_baseline_service_processing_time_synthetic(target, request, random_seed):
    topology = "-".join(request.split("-")[:2])
    num = len(service_reuse[topology])
    random.seed(random_seed)
    random_numbers = [random.gauss(500, 300) for i in range(num)]
    random_numbers = [abs(r) + 200 for r in random_numbers]   # just in case
    random_numbers.sort()
    print(f"[config.py] Random numbers for execution time: {random_numbers}")

    processing_time = {}
    for i in range(num):
        if i in leaf_nodes[topology] and f"service{i}" == target:
            processing_time[f"service{i}"] = round(2 * random_numbers.pop(-1)  / service_reuse[topology][i], 2)
        elif f"service{i}" == target:
            processing_time[f"service{i}"] = round(random_numbers.pop(-1)  / service_reuse[topology][i], 2)
        else:
            processing_time[f"service{i}"] = round(random_numbers.pop(-1)  / service_reuse[topology][i], 2)
    return processing_time

def get_cpu_quota(benchmark, request):
    if benchmark == "synthetic":
        return get_cpu_quota_synthetic(request)
    elif benchmark == "boutique":
        return {
            "cart":2,
            "checkout":2,
            "currency":2,
            "email":2,
            "frontend":2,
            "payment":2,
            "product_catalog":2,
            "recommendations":2,
            "shipping":2
        }
    elif benchmark == "social":
        return {
            "composepost": 2,
            "hometimeline": 2,
            "poststorage": 2,
            "socialgraph": 2,
            "usertimeline": 2
        }
    elif benchmark == "movie":
        return {
            "castinfo": 2,
            "composereview": 2,
            "frontend": 2,
            "movieid": 2,
            "movieinfo": 2,
            "moviereviews": 2,
            "page": 2,
            "plot": 2,
            "reviewstorage": 2,
            "uniqueid": 2,
            "user": 2,
            "userreviews": 2
        }
    elif benchmark == "hotel":
        return {
            "frontend": 2,
            "profile": 2,
            "rate": 2,
            "reservation": 2,
            "search": 2,
            "user": 2
        }
    else :
        raise ValueError(f"[config.py] Unknown benchmark: {benchmark}")

def get_cpu_quota_synthetic(request):
    topology = "-".join(request.split("-")[:2])
    num = len(service_reuse[topology])
    cpu_quota = {}
    for i in range(num):
        cpu_quota[f"service{i}"] = 2
    return cpu_quota