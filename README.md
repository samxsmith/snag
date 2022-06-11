# Snag

The Structured Note Aggregator. Turns your never-to-be-seen-again journals into structured queriable logs, with a flexible schema, creating views, insights & metrics buried between the pages.

Snag is committed to being agnostic to your note editor, file structure, and supports any text-based format, like .txt or .md.

## Problem Statement
I make daily notes, to record what I did each day.
This keeps notes indexed by date, but doesn't allow we to pull by subject or type of activity.

I want to be able to ask questions of my journals and get answers, without having to do extra work outside of my daily notes.

- How many books did read this month or this year?
- What was my favourite book this year?
- How many concerts did I go to this year & what did I think of them?
- What were my favourite restaurants of the year? What did I have?

I want all of this, without removing this information from the context of the day, so I can go back and read by day, and read around these occurrences to see what else was happening.

## Status: beta

I built this project for my own use, and depend upon it daily. Beyond my use-case the project is in beta.
New functionality is coming weekly.

## How to use
My use case is well described by the test example in `e2e/`. I want to structure my notes so I can query and sort them.
For example, books I've read over time can be pulled from my daily journal and sorted by date, or by rating, or both.

- Put a `.snag.yml` at the root of your notes directory.
- Run `snag` in the same directory.
- `snag` will output your aggregated files as you specifcy.

### Config
The `.snag.yml` file has the following options.

```yml
# schemas are the definitions you want to aggregate
schemas:
	# tag is what each entry should start with
	- tag: "book"
	
	# cols are all the columns you'll add to each entry
	  cols:
	  - name
	  - finishDate
	  - rating

	  # sortColumns define how your entries should be sorted
	  # e.g. sort by most recently read, then by rating
	  sortColumns:
	      - name: finishDate
	        type: date
		ascending: false
	      - name: rating
	        type: number

	  # aggregationtype defines how they should be collected
	  # and output. Only list is available at the moment.
	  aggregationType: list

	  # where to save the result
	  outputFilepath: "./aggregated/books_I_read.md"

```

### Entries
With the current config, imagine you've got the following entries in your journals across various days.

> 6th June
> #book The Count of Monte Cristo, 2022-06-06, 5
> Love this book, so glad I re-read....

> 10th May
> #book Walden, 2022-05-10, 4.5

> 2nd Sept
> #book Where the Crawdads Sing, 2021-09-02, 4.5

With structured entries like this and the config file, I'd now get an aggregated list of books I'd read, in a single place.

