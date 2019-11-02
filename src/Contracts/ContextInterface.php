<?php

namespace Lumen\Extra\Contracts;

interface ContextInterface
{
    /**
     * @param string $abstract
     * @param mixed $value
     */
    public function set(string $abstract, $value);

    /**
     * @param $abstract
     * @return mixed
     */
    public function get($abstract);
}
