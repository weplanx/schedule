<?php

namespace lumen\extra\common;

use Illuminate\Pipeline\Pipeline;
use Predis\Client;
use Predis\Transaction\MultiExec;

abstract class RedisModel
{
    /**
     * Model key
     * @var string $key
     */
    protected $key;

    /**
     * Redis Manager
     * @var  Client $redis
     */
    protected $redis;

    /**
     * RedisModel constructor.
     * @param Client|Pipeline|MultiExec $redis
     */
    public function __construct($redis = null)
    {
        $this->redis = $redis ? $redis : app()->make('redis');
    }
}
