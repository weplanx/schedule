<?php

namespace Lumen\Extra\Common;

/**
 * Class ContextFactory
 * @package Lumen\Extra\Common
 */
final class ContextFactory
{
    /**
     * @var array
     */
    private $context = [];

    /**
     * @param string $abstract
     * @param mixed $value
     */
    public function set(string $abstract, $value)
    {
        $this->context[$abstract] = $value;
    }

    /**
     * @param $abstract
     * @return mixed
     */
    public function get($abstract)
    {
        return !empty($this->context[$abstract]) ?
            $this->context[$abstract] :
            null;
    }
}
