let assert = require('chai').assert;

let PriorityQueue_Test = artifacts.require("PriorityQueue_Test");
let { catchError } = require('../utilities.js');

contract('PriorityQueue', async (accounts) => {
    let instance;
    beforeEach(async () => {
        instance = await PriorityQueue_Test.new();
    });

    it("Reverts when deleting on an empty queue", async() => {
      let err;
      [err] = await catchError(instance.getMin.call());
      if (!err)
          assert.fail("Didn't revert on getting min of an empty queue.");

      [err] = await catchError(instance.delMin());
      if (!err)
          assert.fail("Didn't revert on deleting min of an empty queue.");
    });

    it("Correctly adds then remove elements", async () => {
        await instance.insert(2);
        await instance.insert(1);
        await instance.insert(3);

        assert.equal((await instance.getMin.call()).toNumber(), 1, "Did not delete correct minimum");

        await instance.delMin();
        assert.equal((await instance.getMin.call()).toNumber(), 2, "Did not delete correct minimum");

        await instance.delMin();
        assert.equal((await instance.getMin.call()).toNumber(), 3, "Did not delete correct minimum");

        await instance.delMin();
        assert.equal((await instance.currentSize.call()).toNumber(), 0, "Size is not zero");
    });

    it("Handles ascending inserts", async () => {
        let currSize = (await instance.currentSize.call()).toNumber();
        for (i = 1; i < 6; i++) {
            await instance.insert(i);
        }

        currSize = (await instance.currentSize.call()).toNumber();
        assert.equal(currSize, 5, "currentSize did not increment");

        let min = (await instance.getMin.call()).toNumber();
        assert.equal(1, min, "getMin did not return the minimum");

        for (i = 0; i < 3; i++) {
            await instance.delMin();
        }

        min = (await instance.getMin.call()).toNumber();
        currSize = (await instance.currentSize.call()).toNumber();
        assert.equal(min, 4, "delMin deleted priorities out of order");
        assert.equal(currSize, 2, "currSize did not decrement");

        // Clear the queue
        for (i = 0; i < 2; i++) {
            await instance.delMin();
        }
        currSize = (await instance.currentSize.call()).toNumber();
        assert.equal(currSize, 0, "The priority queue has not been emptied");
    });

    it("Can insert, delete min, then insert again", async () => {
        for (i = 1; i < 6; i++) {
            await instance.insert(i);
            let min = (await instance.getMin.call()).toNumber();
            assert.equal(min, 1, "getMin does not return minimum element in pq.");
        }

        // partially clear the pq
        for (i = 0; i < 3; i++) {
            await instance.delMin();
        }
        min = (await instance.getMin.call()).toNumber();
        assert.equal(min, 4, "delMin deleted priorities out of order");

        // insert to pq after partial delete
        for (i = 2; i < 4; i++) {
            await instance.insert(i);
            let min = (await instance.getMin.call()).toNumber();
            assert.equal(min, 2, "getMin does not return minimum element in pq.");
        }
        // clear the pq
        for (i = 0; i < 4; i++) {
            await instance.delMin();
        }

        currSize = (await instance.currentSize.call()).toNumber();
        assert.equal(currSize, 0, "The priority queue has not been emptied");
    });

    it("Handles duplicate entries", async () => {
        let currentSize = (await instance.currentSize.call()).toNumber();
        assert.equal(currentSize, 0, "The size is not 0");

        await instance.insert(10);
        let min = (await instance.getMin.call()).toNumber();
        currentSize = (await instance.currentSize.call()).toNumber();
        assert.equal(currentSize, 1, "The size is not 0");

        // Breaks here - has min as 0
        assert.equal(min, 10, "First insert did not work");

        await instance.insert(10);
        min = (await instance.getMin.call()).toNumber();
        assert.equal(min, 10, "Second insert of same priority did not work");
        await instance.insert(5);
        await instance.insert(5);

        currentSize = (await instance.currentSize.call()).toNumber();
        assert.equal(currentSize, 4, "The currentSize is incorrect")

        await instance.delMin();
        min = (await instance.getMin.call()).toNumber();
        assert.equal(min, 5, "PriorityQueue did not handle same priorities correctly");

        await instance.delMin();
        await instance.delMin();

        min = (await instance.getMin.call()).toNumber();
        assert.equal(min, 10, "PriorityQueue did not delete duplicate correctly")

        await instance.delMin();
        assert.equal(await instance.currentSize.call(), 0);
    });
});
