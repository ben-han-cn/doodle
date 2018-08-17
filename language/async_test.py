import asyncio

async def dd():
    return 20

async def fuck():
    yield 10

async def my_get():
    async for i in fuck():
        print(i)
        a = await dd()
        print(a)


loop = asyncio.get_event_loop()
loop.run_until_complete(my_get())
loop.close()
