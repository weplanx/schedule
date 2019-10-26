<?php

namespace lumen\extra\common;

/**
 * Class RedisModel
 * @package lumen\extra\common
 */
abstract class RedisModel
{
    /**
     * Model key
     * @var string $key
     */
    protected $key;

    /**
     * Redis Manager
     * @var  \Redis $redis
     */
    protected $redis;

    /**
     * Create RedisModel
     * @param \Redis $redis
     * @return static
     */
    public static function create()
    {
        return new static();
    }

    /**
     * RedisModel constructor.
     * @param \Redis $redis
     */
    public function __construct($redis = null)
    {
        $this->redis = $redis ? $redis : app()->make('redis');
    }
}
