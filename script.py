import math
import random

ratings = {
    10: 22,
    9: 60,
    8: 147,
    7: 182,
    6: 130,
    5: 79,
    4: 68,
    3: 29,
    2: 3,
    1: 0
}

total_ratings = sum(ratings.values())
ratings_percentage = []
for i in range(10):
    rating = i + 1
    ratings_percentage.extend([math.pow(2, i)] * ratings[rating])

total = sum(ratings_percentage)
for index, i in enumerate(ratings_percentage):
    ratings_percentage[index] = i/total 

final = [ratings_percentage[0]]
for index, i in enumerate(ratings_percentage[1:]):
    final.append(final[index] + i)

def pick_album(final):
  result = random.random() * final[-1]
  for index, i in enumerate(final):
    if result <= i:
      return index + 1

def get_rating_from_index(index):
    cumulative_ratings = [ratings[1]]
    for i in range(1, 10):
        cumulative_ratings.append(cumulative_ratings[i-1] + ratings[i+1])
    for index, i in enumerate(cumulative_ratings):
        if album_index <= i:
            return index + 1
    return -1


#results = {}
#for i in range(100000):
#  album_index = pick_album(final)
#  rating = get_rating_from_index(album_index)
  #if rating not in results.keys():
      #results[rating] = 0
  #results[rating] += 1

album_index = pick_album(final)
page = (album_index // 25) 
album_no = (album_index - (page * 25)) 
print(album_index, page + 1, album_no + 1)


